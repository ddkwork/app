package translate

import (
	cast "clang/ast"
	"go/ast"
	"go/token"
	"reflect"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

var operatorCodeMap = map[cast.OpCode]token.Token{
	cast.Add: token.ADD, // +
	cast.Sub: token.SUB, // -
	cast.Mul: token.MUL, // *
	cast.Div: token.QUO, // /
	cast.Rem: token.REM, // %

	cast.And: token.AND, // &
	cast.Or:  token.OR,  // |
	cast.Xor: token.XOR, // ^
	cast.Shl: token.SHL, // <<
	cast.Shr: token.SHR, // >>

	cast.LAnd: token.LAND, // &&
	cast.LOr:  token.LOR,  // ||

	cast.EQ: token.EQL, // ==
	cast.NE: token.NEQ, // !=
	cast.LT: token.LSS, // <
	cast.LE: token.LEQ, // <=
	cast.GT: token.GTR, // >
	cast.GE: token.GEQ, // >=
}

var operatorCodeAssignMap = map[cast.OpCode]token.Token{
	cast.Assign: token.ASSIGN, // =

	cast.AddAssign: token.ADD_ASSIGN, // +=
	cast.SubAssign: token.SUB_ASSIGN, // -=
	cast.MulAssign: token.MUL_ASSIGN, // *=
	cast.DivAssign: token.QUO_ASSIGN, // /=
	cast.RemAssign: token.REM_ASSIGN, // %=

	cast.AndAssign: token.AND_ASSIGN, // &=
	cast.OrAssign:  token.OR_ASSIGN,  // |=
	cast.XorAssign: token.XOR_ASSIGN, // ^=
	cast.ShlAssign: token.SHL_ASSIGN, // <<=
	cast.ShrAssign: token.SHR_ASSIGN, // >>=
}

var operatorCodeUnaryMap = map[cast.OpCode]token.Token{
	cast.Add:     token.ADD, // +
	cast.Sub:     token.SUB, // -
	cast.Mul:     token.MUL, // *
	cast.And:     token.AND, // &
	cast.Not:     token.XOR, // ~
	cast.LNot:    token.NOT, // !
	cast.PostInc: token.INC, // ++
	cast.PostDec: token.DEC, // --
}

var (
	thisExpr  = &ast.Ident{Name: "self"}
	nullExpr  = &ast.Ident{Name: "nil"}
	trueExpr  = &ast.Ident{Name: "true"}
	falseExpr = &ast.Ident{Name: "false"}
)

func (ctx *TransCtx) newFuncLitCallExpr(stmts []ast.Stmt, retType string) (ast.Expr, error) {
	t := mylog.Check2(NewTypeParse(retType).Get())

	fnt := &ast.FuncType{Params: &ast.FieldList{}, Results: &ast.FieldList{List: []*ast.Field{{Type: t}}}}
	fn := &ast.FuncLit{Type: fnt, Body: &ast.BlockStmt{List: stmts}}
	call := &ast.CallExpr{Fun: fn}
	return call, nil
}

func (ctx *TransCtx) newConstruct(args []ast.Expr, t string) (ast.Expr, error) {
	tp := mylog.Check2(NewTypeParse(t).Get())

	name, ok := tp.(*ast.Ident)
	if ok {
		name.Name = "New" + name.Name
	}

	if ok && len(args) == 1 {
		if fn, ok := args[0].(*ast.CallExpr); ok {
			if n, ok := fn.Fun.(*ast.Ident); ok && n.Name == name.Name {
				return fn, nil
			}
		}
	}
	fn := &ast.CallExpr{Fun: tp, Args: args}
	return fn, nil
}

func (ctx *TransCtx) transSubValueExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr(node.Kind, " without base value")
	}
	return ctx.transExpr(node.Inner[0])
}

// transArraySubscriptExpr : 数组成员访问
func (ctx *TransCtx) transArraySubscriptExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown array name or index")
	}
	array := mylog.Check2(ctx.transExpr(node.Inner[0]))

	index := mylog.Check2(ctx.transExpr(node.Inner[1]))

	expr := &ast.IndexExpr{X: array, Index: index}
	return expr, nil
}

// transAtomicExpr : 原子操作: __atomic_load_n, __atomic_store_n 等
func (ctx *TransCtx) transAtomicExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("atomic operation no object")
	}
	call := &ast.CallExpr{Fun: &ast.Ident{Name: "atomic"}}
	for _, arg := range node.Inner {
		expr := mylog.Check2(ctx.transExpr(arg))
	}
	return call, nil
}

// transBinaryOperator : 二元运算
func (ctx *TransCtx) transBinaryOperator(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown left or right value for binary operator")
	}
	left := mylog.Check2(ctx.transExpr(node.Inner[0]))

	right := mylog.Check2(ctx.transExpr(node.Inner[1]))

	if op, ok := operatorCodeMap[node.OpCode]; ok {
		expr := &ast.BinaryExpr{X: left, Op: op, Y: right}
		return expr, nil
	}
	if op, ok := operatorCodeAssignMap[node.OpCode]; ok {
		if node.Type == nil {
			return nil, ctx.newErr("unknown express type")
		}
		stmt := &ast.AssignStmt{Lhs: []ast.Expr{left}, Tok: op, Rhs: []ast.Expr{right}}
		retStmt := &ast.ReturnStmt{Results: []ast.Expr{left}}
		return ctx.newFuncLitCallExpr([]ast.Stmt{stmt, retStmt}, node.Type.QualType)
	}

	return nil, ctx.newErr("unsupport op code: ", node.OpCode)
}

// transCXXBoolLiteralExpr : true, false
func (ctx *TransCtx) transCXXBoolLiteralExpr(node *cast.Node) (ast.Expr, error) {
	val, ok := node.Value.(bool)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Value)
	}
	if val {
		return trueExpr, nil
	} else {
		return falseExpr, nil
	}
}

// transCXXConstructExpr : C++变量构造
func (ctx *TransCtx) transCXXConstructExpr(node *cast.Node) (ast.Expr, error) {
	if node.Type == nil {
		return nil, ctx.newErr("constant without base value and type")
	}
	var args []ast.Expr
	for _, parm := range node.Inner {
		if parm.Kind == cast.CXXDefaultArgExpr {
			continue
		}
		expr := mylog.Check2(ctx.transExpr(parm))

	}
	return ctx.newConstruct(args, node.Type.QualType)
}

// transCXXDeleteExpr : C++ delete表达式
func (ctx *TransCtx) transCXXDeleteExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("no object to delete")
	}

	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	fn := &ast.CallExpr{Fun: &ast.Ident{Name: "delete"}, Args: []ast.Expr{expr}}
	return fn, nil
}

// transCXXDependentScopeMemberExpr : C++虚成员
func (ctx *TransCtx) transCXXDependentScopeMemberExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("no function to call")
	}

	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	name := transName(node.Member)
	fn := &ast.SelectorExpr{X: expr, Sel: &ast.Ident{Name: name}}
	return fn, nil
}

// transCXXMemberCallExpr : C++方法调用
func (ctx *TransCtx) transCXXMemberCallExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("no function to call")
	}

	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	fn := &ast.CallExpr{Fun: expr}
	for _, parm := range node.Inner[1:] {
		if parm.Kind == cast.CXXDefaultArgExpr {
			continue
		}
		expr = mylog.Check2(ctx.transExpr(parm))

	}
	return fn, nil
}

// transCXXOperatorCallExpr : C++运算符重载方法调用
func (ctx *TransCtx) transCXXOperatorCallExpr(node *cast.Node) (ast.Expr, error) {
	var arg []ast.Expr
	for _, n := range node.Inner {
		expr := mylog.Check2(ctx.transExpr(n))
	}
	if len(arg) < 2 {
		return nil, ctx.newErr("no function to call")
	}

	name, ok := arg[0].(*ast.Ident)
	if !ok {
		if sel, ok := arg[0].(*ast.SelectorExpr); ok {
			name = sel.Sel
		} else {
			return nil, ctx.newErr("unknow function name: ", reflect.TypeOf(arg[0]))
		}
	}
	name.Name = transName(name.Name)

	fn := &ast.CallExpr{Fun: &ast.SelectorExpr{X: arg[1], Sel: name}, Args: arg[2:]}
	return fn, nil
}

// transCXXTemporaryObjectExpr : C++临时的值: return String("hello")
func (ctx *TransCtx) transCXXTemporaryObjectExpr(node *cast.Node) (ast.Expr, error) {
	if node.Type == nil {
		return nil, ctx.newErr("unknown object type")
	}

	var args []ast.Expr
	for _, parm := range node.Inner {
		expr := mylog.Check2(ctx.transExpr(parm))
	}
	return ctx.newConstruct(args, node.Type.QualType)
}

// transCXXThrowExpr : C++ throw表达式
func (ctx *TransCtx) transCXXThrowExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("nothing to throw")
	}

	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	expr = &ast.CallExpr{Fun: &ast.Ident{Name: "panic"}, Args: []ast.Expr{expr}}
	return expr, nil
}

// transCXXCtorInitializer : C++构造函数初始化: ptr(NULL)
func (ctx *TransCtx) transCXXCtorInitializer(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown value to initial")
	}

	name := ""
	var val ast.Expr

	if node.AnyInit != nil {
		name = node.AnyInit.Name
		val = mylog.Check2(ctx.transExpr(node.Inner[0]))

	} else if node.BaseInit != nil {
		name = node.BaseInit.QualType
		var args []ast.Expr
		if node.Inner[0].Kind == cast.ParenListExpr {
			ctx.InNode(node.Inner[0])
			args = ctx.transParenListExpr(node.Inner[0])
		} else {
			arg := mylog.Check2(ctx.transExpr(node.Inner[0]))
		}
		val = mylog.Check2(ctx.newConstruct(args, name))

	} else {
		return nil, ctx.newErr("unknown field to initial")
	}

	expr := &ast.KeyValueExpr{Key: &ast.Ident{Name: transName(name)}, Value: val}
	return expr, nil
}

// transCallExpr : 函数调用
func (ctx *TransCtx) transCallExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("no function to call")
	}
	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	fn := &ast.CallExpr{Fun: expr}
	for _, parm := range node.Inner[1:] {
		expr = mylog.Check2(ctx.transExpr(parm))
	}
	return fn, nil
}

// transCharacterLiteral : 字符常量
func (ctx *TransCtx) transCharacterLiteral(node *cast.Node) (ast.Expr, error) {
	val, ok := node.Value.(float64)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Value)
	}
	str := string([]byte{'\'', byte(val), '\''})
	switch byte(val) {
	case '\'':
		str = "'\\''"
	case '\\':
		str = "'\\\\'"
	case '\r':
		str = "'\\r'"
	case '\n':
		str = "'\\n'"
	case '\v':
		str = "'\\v'"
	case '\f':
		str = "'\\f'"
	case 0:
		str = "0"
	default:
		break
	}
	expr := &ast.BasicLit{Kind: token.CHAR, Value: str}
	return expr, nil
}

// transCompoundLiteralExpr : 常量结构体/数组: (int[]){1, 2}
func (ctx *TransCtx) transCompoundLiteralExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("compound literal without base value")
	}
	// 包含 InitListExpr 子结构
	return ctx.transExpr(node.Inner[0])
}

// transConditionalOperator : 三元运算: (cond ? expr1 : expr2)
func (ctx *TransCtx) transConditionalOperator(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) < 3 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[0]))

	// 成功分支
	expr1 := mylog.Check2(ctx.transExpr(node.Inner[1]))

	// 失败分支
	expr2 := mylog.Check2(ctx.transExpr(node.Inner[2]))

	if node.Type == nil {
		return nil, ctx.newErr("unknown express type")
	}
	stmt := &ast.IfStmt{Cond: cond}
	stmt.Body = &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{expr1}}}}
	stmt.Else = &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{expr2}}}}
	return ctx.newFuncLitCallExpr([]ast.Stmt{stmt}, node.Type.QualType)
}

// transConstantExpr : 常量
func (ctx *TransCtx) transConstantExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("constant without base value")
	}
	return ctx.transExpr(node.Inner[0])
}

// transCStyleCastExpr : C风格类型转换
func (ctx *TransCtx) transCStyleCastExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("cast without base value")
	}

	if node.CastKind == cast.NullToPointer && node.Inner[0].Value == "0" {
		return nullExpr, nil
	}
	arg := mylog.Check2(ctx.transExpr(node.Inner[0]))

	expr := &ast.CallExpr{Args: []ast.Expr{arg}}
	expr.Fun = mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	return expr, nil
}

// transDeclRefExpr : 变量
func (ctx *TransCtx) transDeclRefExpr(node *cast.Node) (ast.Expr, error) {
	if node.ReferencedDecl == nil || node.ReferencedDecl.Name == "" {
		return nil, ctx.newErr("referenced without name")
	}
	expr := &ast.Ident{Name: transName(node.ReferencedDecl.Name)}
	if node.ReferencedDecl.Kind == cast.CXXMethodDecl {
		return &ast.SelectorExpr{X: thisExpr, Sel: expr}, nil
	}
	return expr, nil
}

// transFloatingLiteral : 浮点数常量
func (ctx *TransCtx) transFloatingLiteral(node *cast.Node) (ast.Expr, error) {
	val, ok := node.Value.(string)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Value)
	}
	expr := &ast.BasicLit{Kind: token.FLOAT, Value: val}
	return expr, nil
}

// transImaginaryLiteral : 虚数常量
func (ctx *TransCtx) transImaginaryLiteral(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("imaginary without base value")
	}
	ctx.InNode(node.Inner[0])
	if node.Inner[0].Kind != cast.IntegerLiteral && node.Inner[0].Kind != cast.FloatingLiteral {
		return nil, ctx.newKindErr(node.Inner[0])
	}
	val, ok := node.Inner[0].Value.(string)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Inner[0].Value)
	}
	expr := &ast.BasicLit{Kind: token.IMAG, Value: val + "i"}
	return expr, nil
}

// transImplicitCastExpr : 隐式类型转换
func (ctx *TransCtx) transImplicitCastExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("cast without base value")
	}
	arg := mylog.Check2(ctx.transExpr(node.Inner[0]))

	switch node.CastKind {
	case cast.FunctionToPointerDecay:
		// 函数指针转换
		return arg, nil
	case cast.LValueToRValue:
		// 左右值转换
		return arg, nil
	case cast.NullToPointer:
		return nullExpr, nil
	default:
		break
	}
	if 0 < 1 {
		return arg, nil
	}
	if _, ok := arg.(*ast.BasicLit); ok {
		// 常数
		return arg, nil
	}

	if node.Type == nil {
		return nil, ctx.newErr("cast without target type")
	}
	expr := &ast.CallExpr{Args: []ast.Expr{arg}}
	expr.Fun = mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	return expr, nil
}

// transImplicitValueInitExpr : 结构体隐式初始化
func (ctx *TransCtx) transImplicitValueInitExpr(node *cast.Node) (ast.Expr, error) {
	if node.Type == nil {
		return nil, ctx.newErr("init value without type")
	}
	t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	if tp, ok := t.(*ast.Ident); ok {
		switch tp.Name {
		case "byte", "int8", "uint8", "int16", "uint16", "int32", "uint32",
			"int64", "uint64", "int", "uint", "uintptr", "rune":
			return &ast.BasicLit{Kind: token.INT, Value: "0"}, nil
		case "float32", "float64", "complex64", "complex128":
			return &ast.BasicLit{Kind: token.FLOAT, Value: "0.0"}, nil
		default:
			return &ast.CompositeLit{Type: t}, nil
		}
	}
	return &ast.Ident{Name: "nil"}, nil
}

// transInitListExpr : 初始化结构列表: {1, 2}
func (ctx *TransCtx) transInitListExpr(node *cast.Node) (ast.Expr, error) {
	if node.Type == nil {
		return nil, ctx.newErr("init list without type")
	}
	t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	expr := &ast.CompositeLit{Type: t}
	for _, e := range node.Inner {
		elem := mylog.Check2(ctx.transExpr(e))
	}
	return expr, nil
}

// transIntegerLiteral : 整数常量
func (ctx *TransCtx) transIntegerLiteral(node *cast.Node) (ast.Expr, error) {
	val, ok := node.Value.(string)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Value)
	}
	expr := &ast.BasicLit{Kind: token.INT, Value: val}
	return expr, nil
}

// transMemberExpr : 结构体成员访问
func (ctx *TransCtx) transMemberExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("member without base value")
	}
	base := mylog.Check2(ctx.transExpr(node.Inner[0]))

	expr := &ast.SelectorExpr{X: base, Sel: &ast.Ident{Name: transName(node.Name)}}
	return expr, nil
}

// transOffsetOfExpr : offsetof 运算
func (ctx *TransCtx) transOffsetOfExpr(node *cast.Node) (ast.Expr, error) {
	expr := &ast.CallExpr{Fun: &ast.Ident{Name: "offsetof"}}
	return expr, nil
}

// transParenExpr : 括号
func (ctx *TransCtx) transParenExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("paren without base value")
	}
	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	switch expr.(type) {
	case *ast.Ident, *ast.BasicLit:
		return expr, nil
	default:
		break
	}
	expr = &ast.ParenExpr{X: expr}
	return expr, nil
}

// transParenListExpr : 括号列表, 用于C++构造函数初始化: (1, 2)
func (ctx *TransCtx) transParenListExpr(node *cast.Node) []ast.Expr {
	var exprs []ast.Expr
	for _, n := range node.Inner {
		expr := mylog.Check2(ctx.transExpr(n))
	}
	return exprs
}

// transPredefinedExpr : 预定义常量: __func__, __LINE__, __FILE__ 等
func (ctx *TransCtx) transPredefinedExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("predefined without real value")
	}
	utils.Log(utils.LL_Debug, "predefined %s load", node.Name)
	return ctx.transExpr(node.Inner[0])
}

// transRecoveryExpr : 接收器
func (ctx *TransCtx) transRecoveryExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		utils.Log(utils.LL_Warn, "RecoveryExpr without base value")
		return nullExpr, nil
	}
	return ctx.transExpr(node.Inner[0])
}

// transStringLiteral : 字符串常量
func (ctx *TransCtx) transStringLiteral(node *cast.Node) (ast.Expr, error) {
	val, ok := node.Value.(string)
	if !ok {
		return nil, ctx.newErr("unknown value type: ", node.Value)
	}
	expr := &ast.BasicLit{Kind: token.STRING, Value: val}
	return expr, nil
}

// transUnaryExprOrTypeTraitExpr : sizeof 运算
func (ctx *TransCtx) transUnaryExprOrTypeTraitExpr(node *cast.Node) (ast.Expr, error) {
	expr := &ast.CallExpr{Fun: &ast.Ident{Name: node.Name}}
	if len(node.Inner) > 0 {
		arg := mylog.Check2(ctx.transExpr(node.Inner[0]))
	} else if node.ArgType != nil {
		t := mylog.Check2(NewTypeParse(node.ArgType.QualType).Get())
	}
	return expr, nil
}

// transUnaryOperator : 一元运算
func (ctx *TransCtx) transUnaryOperator(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unary operator without base value")
	}
	expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

	op, ok := operatorCodeUnaryMap[node.OpCode]
	if !ok {
		utils.Err(ctx.newErr("unsupport op code:", node.OpCode))
	}
	unary := &ast.UnaryExpr{Op: op, X: expr}
	return unary, nil
}

// transUnresolvedLookupExpr : 未定义符号
func (ctx *TransCtx) transUnresolvedLookupExpr(node *cast.Node) (ast.Expr, error) {
	if node.Name != "" {
		utils.Log(utils.LL_Warn, "undeclared identifier '%s'", node.Name)
		return &ast.Ident{Name: transName(node.Name)}, nil
	}
	if len(node.Inner) == 0 {
		return nullExpr, nil
	}
	return ctx.transExpr(node.Inner[0])
}

// transVAArgExpr : va_arg 运算
func (ctx *TransCtx) transVAArgExpr(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("va_arg() without base value")
	}
	arg := mylog.Check2(ctx.transExpr(node.Inner[0]))

	if node.Type == nil {
		return nil, ctx.newErr("va_arg() without value type")
	}
	t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	expr := &ast.CallExpr{Fun: &ast.Ident{Name: "va_arg"}, Args: []ast.Expr{arg, t}}
	return expr, nil
}

func (ctx *TransCtx) transExpr(node *cast.Node) (ast.Expr, error) {
	ctx.InNode(node)

	switch node.Kind {
	case cast.ArraySubscriptExpr:
		return ctx.transArraySubscriptExpr(node)
	case cast.AtomicExpr:
		return ctx.transAtomicExpr(node)
	case cast.BinaryOperator:
		return ctx.transBinaryOperator(node)
	case cast.CXXBoolLiteralExpr:
		return ctx.transCXXBoolLiteralExpr(node)
	case cast.CXXBindTemporaryExpr:
		return ctx.transSubValueExpr(node)
	case cast.CXXConstructExpr:
		return ctx.transCXXConstructExpr(node)
	case cast.CXXDeleteExpr:
		return ctx.transCXXDeleteExpr(node)
	case cast.CXXDependentScopeMemberExpr:
		return ctx.transCXXDependentScopeMemberExpr(node)
	case cast.CXXFunctionalCastExpr:
		return ctx.transSubValueExpr(node)
	case cast.CXXMemberCallExpr:
		return ctx.transCXXMemberCallExpr(node)
	case cast.CXXNewExpr:
		return ctx.transSubValueExpr(node)
	case cast.CXXOperatorCallExpr:
		return ctx.transCXXOperatorCallExpr(node)
	case cast.CXXTemporaryObjectExpr:
		return ctx.transCXXTemporaryObjectExpr(node)
	case cast.CXXThisExpr:
		return thisExpr, nil
	case cast.CXXThrowExpr:
		return ctx.transCXXThrowExpr(node)
	case cast.CXXUnresolvedConstructExpr:
		return ctx.transSubValueExpr(node)
	case cast.CallExpr:
		return ctx.transCallExpr(node)
	case cast.CharacterLiteral:
		return ctx.transCharacterLiteral(node)
	case cast.CompoundLiteralExpr:
		return ctx.transCompoundLiteralExpr(node)
	case cast.ConditionalOperator:
		return ctx.transConditionalOperator(node)
	case cast.ConstantExpr:
		return ctx.transConstantExpr(node)
	case cast.CStyleCastExpr:
		return ctx.transCStyleCastExpr(node)
	case cast.DeclRefExpr:
		return ctx.transDeclRefExpr(node)
	case cast.ExprWithCleanups:
		return ctx.transSubValueExpr(node)
	case cast.FloatingLiteral:
		return ctx.transFloatingLiteral(node)
	case cast.ImaginaryLiteral:
		return ctx.transImaginaryLiteral(node)
	case cast.ImplicitCastExpr:
		return ctx.transImplicitCastExpr(node)
	case cast.ImplicitValueInitExpr:
		return ctx.transImplicitValueInitExpr(node)
	case cast.InitListExpr:
		return ctx.transInitListExpr(node)
	case cast.IntegerLiteral:
		return ctx.transIntegerLiteral(node)
	case cast.MaterializeTemporaryExpr:
		return ctx.transSubValueExpr(node)
	case cast.MemberExpr:
		return ctx.transMemberExpr(node)
	case cast.OffsetOfExpr:
		return ctx.transOffsetOfExpr(node)
	case cast.ParenExpr:
		return ctx.transParenExpr(node)
	case cast.PredefinedExpr:
		return ctx.transPredefinedExpr(node)
	case cast.RecoveryExpr:
		return ctx.transRecoveryExpr(node)
	case cast.StringLiteral:
		return ctx.transStringLiteral(node)
	case cast.UnaryExprOrTypeTraitExpr:
		return ctx.transUnaryExprOrTypeTraitExpr(node)
	case cast.UnaryOperator:
		return ctx.transUnaryOperator(node)
	case cast.UnresolvedLookupExpr:
		return ctx.transUnresolvedLookupExpr(node)
	case cast.VAArgExpr:
		return ctx.transVAArgExpr(node)
	default:
		return nil, ctx.newKindErr(node)
	}
}
