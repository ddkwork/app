package translate

import (
	cast "clang/ast"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

var operatorNameMap = map[token.Token]string{
	token.ADD: "Add", // +
	token.SUB: "Sub", // -
	token.MUL: "Mul", // *
	token.QUO: "Div", // /
	token.REM: "Rem", // %

	token.AND: "And",        // &
	token.OR:  "Or",         // |
	token.XOR: "Xor",        // ^
	token.SHL: "ShiftLeft",  // <<
	token.SHR: "ShiftRight", // >>

	token.LAND: "LAnd", // &&
	token.LOR:  "LOr",  // ||
	token.NOT:  "LNot", // !

	token.EQL: "Equal",        // ==
	token.NEQ: "NotEqual",     // !=
	token.LSS: "Less",         // <
	token.LEQ: "LessEqual",    // <=
	token.GTR: "Greater",      // >
	token.GEQ: "GreaterEqual", // >=

	token.ASSIGN: "Assign", // =

	token.ADD_ASSIGN: "AddAssign", // +=
	token.SUB_ASSIGN: "SubAssign", // -=
	token.MUL_ASSIGN: "MulAssign", // *=
	token.QUO_ASSIGN: "DivAssign", // /=
	token.REM_ASSIGN: "RemAssign", // %=

	token.AND_ASSIGN: "AndAssign", // &=
	token.OR_ASSIGN:  "OrAssign",  // |=
	token.XOR_ASSIGN: "XorAssign", // ^=
	token.SHL_ASSIGN: "ShlAssign", // <<=
	token.SHR_ASSIGN: "ShrAssign", // >>=

	token.INC: "PostInc", // ++
	token.DEC: "PostDec", // --
}

func transName(name string) string {
	if name == "" {
		return ""
	}
	switch name {
	// go关键字、内置函数
	case "append", "cap", "chan", "close", "complex", "copy", "defer", "delete", "fallthrough", "func", "go",
		"imag", "import", "interface", "iota", "len", "make", "map", "new", "nil",
		"package", "panic", "range", "real", "recover", "select", "type", "var":
		return name + "_"
	default:
		break
	}
	if strings.HasPrefix(name, "__go_") {
		return name[5:]
	}
	if !strings.HasPrefix(name, "operator") {
		return name
	}

	scan := NewTokenScanner(name)
	_, tok, lit := scan.Scan()
	if tok != token.IDENT || lit != "operator" {
		return name
	}

	_, tok, lit = scan.Scan()
	if str, ok := operatorNameMap[tok]; ok {
		return str
	}
	return name
}

// transCXXRecordDecl : class ClassName { };
func (ctx *TransCtx) transCXXRecordDecl(node *cast.Node) {
	if node.IsImplicit {
		return
	}
	if len(node.Inner) == 0 {
		// 忽略类的声明
		return
	}

	t := &ast.StructType{Fields: &ast.FieldList{}}
	for _, base := range node.Bases {
		if base.Type == nil {
			utils.Err(ctx.newErr("base without type"))
			continue
		}
		bt := mylog.Check2(NewTypeParse(base.Type.QualType).Get())

		field := &ast.Field{Type: &ast.StarExpr{X: bt}}
		t.Fields.List = append(t.Fields.List, field)
	}

	self := &ast.Field{Names: []*ast.Ident{thisExpr}, Type: &ast.StarExpr{X: &ast.Ident{Name: node.Name}}}
	var decls []ast.Decl
	var doc *ast.CommentGroup
	for _, field := range node.Inner {
		ctx.InNode(field)
		switch field.Kind {
		case cast.CXXMethodDecl:
			decl := mylog.Check2(ctx.transCXXMethodDecl(field))

		case cast.FieldDecl:
			f := mylog.Check2(ctx.transFieldDecl(field))

		case cast.FullComment:
			doc = ctx.transFullComment(field, node.Name)
		default:
			decl := mylog.Check2(ctx.transDecl(field))

		}
	}

	spec := &ast.TypeSpec{Name: &ast.Ident{Name: node.Name}, Type: t}
	decl := &ast.GenDecl{Doc: doc, Tok: token.TYPE, Specs: []ast.Spec{spec}}
	ctx.goFile.Decls = append(ctx.goFile.Decls, decl)
	ctx.goFile.Decls = append(ctx.goFile.Decls, decls...)
}

// transCXXConstructorDecl : 构造函数
func (ctx *TransCtx) transCXXConstructorDecl(node *cast.Node) (ast.Decl, error) {
	if node.IsImplicit {
		return nil, nil
	}
	if len(node.Inner) == 0 {
		// 忽略构造函数声明
		return nil, nil
	}

	var params []*ast.Field
	var body *ast.BlockStmt
	var list []ast.Expr
	var doc *ast.CommentGroup
	for _, parm := range node.Inner {
		ctx.InNode(parm)
		switch parm.Kind {
		case cast.ParmVarDecl:
			field := mylog.Check2(ctx.transParmVarDecl(parm))

		case cast.CompoundStmt:
			body = ctx.transCompoundStmt(parm)
		case cast.CXXCtorInitializer:
			init := mylog.Check2(ctx.transCXXCtorInitializer(parm))

		case cast.FullComment:
			doc = ctx.transFullComment(parm, "New"+node.Name)
		default:
			utils.Err(ctx.newKindErr(parm))
		}
	}

	if len(list) > 0 {
		if body == nil {
			body = &ast.BlockStmt{}
		}
	} else if body == nil {
		// 忽略构造函数声明
		return nil, nil
	}

	name := &ast.Ident{Name: node.Name}
	expr := &ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: name, Elts: list}}
	self := &ast.AssignStmt{Lhs: []ast.Expr{thisExpr}, Tok: token.ASSIGN, Rhs: []ast.Expr{expr}}
	ret := &ast.Field{Names: []*ast.Ident{thisExpr}, Type: &ast.StarExpr{X: name}}

	temp := body.List
	body.List = make([]ast.Stmt, 0, len(list)+2)
	body.List = append(body.List, self)
	body.List = append(body.List, temp...)
	body.List = append(body.List, &ast.ReturnStmt{})

	ft := &ast.FuncType{Params: &ast.FieldList{List: params}, Results: &ast.FieldList{List: []*ast.Field{ret}}}
	decl := &ast.FuncDecl{Doc: doc, Name: &ast.Ident{Name: "New" + node.Name}, Type: ft, Body: body}
	return decl, nil
}

// transCXXDestructorDecl : 析构函数
func (ctx *TransCtx) transCXXDestructorDecl(node *cast.Node) (ast.Decl, error) {
	if node.IsImplicit {
		return nil, nil
	}
	if len(node.Inner) == 0 {
		// 忽略析构函数声明
		return nil, nil
	}

	var body *ast.BlockStmt
	var doc *ast.CommentGroup
	for _, parm := range node.Inner {
		ctx.InNode(parm)
		switch parm.Kind {
		case cast.CompoundStmt:
			body = ctx.transCompoundStmt(parm)
		case cast.FullComment:
			doc = ctx.transFullComment(parm, "Close")
		default:
			utils.Err(ctx.newKindErr(parm))
		}
	}

	if body == nil {
		// 忽略析构函数声明
		return nil, nil
	}

	self := &ast.Field{Names: []*ast.Ident{thisExpr}, Type: &ast.StarExpr{X: &ast.Ident{Name: node.Name[1:]}}}
	recv := &ast.FieldList{List: []*ast.Field{self}}
	ft := &ast.FuncType{Params: &ast.FieldList{}}
	decl := &ast.FuncDecl{Doc: doc, Recv: recv, Name: &ast.Ident{Name: "Close"}, Type: ft, Body: body}
	return decl, nil
}

// transCXXMethodDecl : 类的方法
func (ctx *TransCtx) transCXXMethodDecl(node *cast.Node) (ast.Decl, error) {
	if node.IsImplicit {
		return nil, nil
	}
	if node.Type == nil {
		return nil, ctx.newErr("unknown function type")
	}
	if len(node.Inner) == 0 {
		// 忽略方法声明
		return nil, nil
	}

	name := transName(node.Name)
	var params []*ast.Field
	var body *ast.BlockStmt
	var doc *ast.CommentGroup
	for _, parm := range node.Inner {
		ctx.InNode(parm)
		switch parm.Kind {
		case cast.ParmVarDecl:
			field := mylog.Check2(ctx.transParmVarDecl(parm))

		case cast.CompoundStmt:
			body = ctx.transCompoundStmt(parm)
		case cast.FullComment:
			doc = ctx.transFullComment(parm, name)
		default:
			utils.Err(ctx.newKindErr(parm))
		}
	}

	if body == nil {
		// 忽略构造函数声明
		return nil, nil
	}

	parent := ctx.nodes[node.ParentDeclContextID.ToInt()].node
	ret := mylog.Check2(NewTypeParse(node.Type.QualType).GetFuncRet())

	ft := &ast.FuncType{Params: &ast.FieldList{List: params}, Results: ret}
	decl := &ast.FuncDecl{Doc: doc, Name: &ast.Ident{Name: name}, Type: ft, Body: body}
	if parent != nil {
		self := &ast.Field{Names: []*ast.Ident{thisExpr}, Type: &ast.StarExpr{X: &ast.Ident{Name: parent.Name}}}
		decl.Recv = &ast.FieldList{List: []*ast.Field{self}}
	}
	return decl, nil
}

func (ctx *TransCtx) transEnumConstantDecl(node *cast.Node, index int, val *int) *ast.ValueSpec {
	spec := &ast.ValueSpec{Names: []*ast.Ident{{Name: node.Name}}}
	*val++
	if index == 0 {
		spec.Values = []ast.Expr{&ast.Ident{Name: "iota"}}
	}

	if len(node.Inner) == 0 {
		return spec
	}
	var valueNode *cast.Node
	for _, n := range node.Inner {
		ctx.InNode(n)
		switch n.Kind {
		case cast.ConstantExpr:
			valueNode = n
		case cast.FullComment:
			spec.Doc = ctx.transFullComment(n, node.Name)
		default:
			utils.Err(ctx.newKindErr(n))
		}
	}
	if valueNode == nil {
		return spec
	}

	valStr, ok := node.Inner[0].Value.(string)
	if !ok {
		utils.Err(ctx.newErr("unknown value type: ", node.Inner[0].Value))
		return spec
	}
	value := mylog.Check2(strconv.Atoi(valStr))

	if *val != value {
		expr := &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(value - index)}
		spec.Values = []ast.Expr{&ast.BinaryExpr{X: expr, Op: token.ADD, Y: &ast.Ident{Name: "iota"}}}
		*val = value
	}
	return spec
}

func (ctx *TransCtx) transEnumDecl(node *cast.Node) (ast.Decl, error) {
	decl := &ast.GenDecl{Tok: token.CONST, Specs: make([]ast.Spec, 0, len(node.Inner))}
	index, val := 0, -1
	for _, e := range node.Inner {
		ctx.InNode(e)
		switch e.Kind {
		case cast.EnumConstantDecl:
			spec := ctx.transEnumConstantDecl(e, index, &val)
			if len(spec.Values) > 0 && node.Name != "" {
				spec.Type = &ast.Ident{Name: node.Name}
			}
			index++
			decl.Specs = append(decl.Specs, spec)
		case cast.FullComment:
			decl.Doc = ctx.transFullComment(e, node.Name)
		default:
			utils.Err(ctx.newKindErr(e))
		}
	}

	if node.Name != "" {
		spec := &ast.TypeSpec{Name: &ast.Ident{Name: node.Name}, Type: &ast.Ident{Name: "int"}}
		typeDecl := &ast.GenDecl{Doc: decl.Doc, Tok: token.TYPE, Specs: []ast.Spec{spec}}
		ctx.goFile.Decls = append(ctx.goFile.Decls, typeDecl)
		decl.Doc = nil
	}

	ctx.setNodeData(node.ID.ToInt(), decl)
	return decl, nil
}

func (ctx *TransCtx) transFieldDecl(node *cast.Node) (*ast.Field, error) {
	id := node.ID.ToInt()
	if id > 0 && ctx.nodes[id].data != nil {
		if f, ok := ctx.nodes[id].data.(*ast.Field); ok {
			return f, nil
		}
		utils.Log(utils.LL_Error, "node data unexpect: %+v", ctx.nodes[id].data)
	}

	if node.Type == nil {
		return nil, ctx.newErr("unknown field type")
	}
	t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	field := &ast.Field{Names: []*ast.Ident{{Name: transName(node.Name)}}, Type: t}
	if node.Name == "" {
		field.Names[0].Name = "_"
	}

	if len(node.Inner) > 0 && node.Inner[0].Kind == cast.FullComment {
		ctx.InNode(node.Inner[0])
		field.Doc = ctx.transFullComment(node.Inner[0], node.Name)
	}

	ctx.setNodeData(id, field)
	return field, err
}

func (ctx *TransCtx) transFunctionDecl(node *cast.Node) (ast.Decl, error) {
	if len(node.Inner) == 0 {
		// 忽略函数声明
		return nil, nil
	}

	if node.Type == nil {
		return nil, ctx.newErr("unknown function type")
	}
	ret := mylog.Check2(NewTypeParse(node.Type.QualType).GetFuncRet())

	name := transName(node.Name)
	t := &ast.FuncType{Params: &ast.FieldList{}, Results: ret}
	decl := &ast.FuncDecl{Name: &ast.Ident{Name: name}, Type: t}

	for _, parm := range node.Inner {
		ctx.InNode(parm)
		switch parm.Kind {
		case cast.ParmVarDecl:
			field := mylog.Check2(ctx.transParmVarDecl(parm))

		case cast.CompoundStmt:
			decl.Body = ctx.transCompoundStmt(parm)
		case cast.FullComment:
			decl.Doc = ctx.transFullComment(parm, name)
			ctx.setNodeData(node.ID.ToInt(), decl.Doc)
		default:
			utils.Err(ctx.newKindErr(parm))
		}
	}
	if node.StorageClass == "extern" || decl.Body == nil {
		// 忽略函数声明
		return nil, nil
	}

	// 查找在声明处的注释
	previousDecl := node.PreviousDecl
	for previousDecl != "" && decl.Doc == nil {
		obj := ctx.nodes[previousDecl.ToInt()]
		if obj.node == nil {
			break
		}
		decl.Doc, _ = obj.data.(*ast.CommentGroup)
		previousDecl = obj.node.PreviousDecl
	}

	return decl, nil
}

// transLinkageSpecDecl : extren "C" { }
func (ctx *TransCtx) transLinkageSpecDecl(node *cast.Node) {
	for _, n := range node.Inner {
		decl := mylog.Check2(ctx.transDecl(n))
	}
}

// transNamespaceDecl : namespace xxx { }
func (ctx *TransCtx) transNamespaceDecl(node *cast.Node) {
	for _, n := range node.Inner {
		decl := mylog.Check2(ctx.transDecl(n))
	}
}

func (ctx *TransCtx) transParmVarDecl(node *cast.Node) (*ast.Field, error) {
	if node.Type == nil {
		return nil, ctx.newErr("unknown input parameter type")
	}
	t0 := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	field := &ast.Field{Names: []*ast.Ident{{Name: transName(node.Name)}}, Type: t0}
	if node.Name == "" {
		field.Names[0].Name = "_"
	}
	return field, nil
}

func (ctx *TransCtx) transRecordDecl(node *cast.Node) (ast.Decl, error) {
	t := &ast.StructType{Fields: &ast.FieldList{}}
	if len(node.Inner) > 0 {
		t.Fields.List = make([]*ast.Field, 0, len(node.Inner))
		for _, field := range node.Inner {
			ctx.InNode(field)
			if field.Kind != cast.FieldDecl {
				utils.Err(ctx.newKindErr(field))
				continue
			}
			f := mylog.Check2(ctx.transFieldDecl(field))

			t.Fields.List = append(t.Fields.List, f)
		}
	}

	ctx.setNodeData(node.ID.ToInt(), ast.Expr(t))

	if node.Name != "" {
		spec := &ast.TypeSpec{Name: &ast.Ident{Name: node.Name}, Type: t}
		decl := &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{spec}}
		return decl, nil
	}
	return nil, nil
}

func (ctx *TransCtx) transTypedefDecl(node *cast.Node) (ast.Decl, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("type defined without base type")
	}
	t := mylog.Check2(ctx.transType(node.Inner[0]))

	// 重写枚举值类型
	if node.Inner[0].Kind == cast.ElaboratedType && node.Inner[0].OwnedTagDecl != nil &&
		node.Inner[0].OwnedTagDecl.Kind == cast.EnumDecl && node.Inner[0].OwnedTagDecl.Name == "" {
		val, ok := ctx.nodes[node.Inner[0].OwnedTagDecl.ID.ToInt()].data.(*ast.GenDecl)
		if ok && val != nil && val.Tok == token.CONST {
			for _, e := range val.Specs {
				if ev, ok := e.(*ast.ValueSpec); ok {
					ev.Type = &ast.Ident{Name: node.Name}
				}
			}
		}
	}

	var doc *ast.CommentGroup
	if len(node.Inner) > 1 && node.Inner[1].Kind == cast.FullComment {
		doc = ctx.transFullComment(node.Inner[1], node.Name)
	}

	spec := &ast.TypeSpec{Name: &ast.Ident{Name: node.Name}, Type: t}
	decl := &ast.GenDecl{Doc: doc, Tok: token.TYPE, Specs: []ast.Spec{spec}}
	return decl, nil
}

func (ctx *TransCtx) transVarDecl(node *cast.Node) (ast.Decl, error) {
	if node.StorageClass == "extern" {
		// 忽略变量声明
		return nil, nil
	}
	if node.Type == nil {
		return nil, ctx.newErr("var without type")
	}
	t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

	spec := &ast.ValueSpec{Names: []*ast.Ident{{Name: transName(node.Name)}}, Type: t}
	decl := &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{spec}}

	if node.Init != "" && len(node.Inner) > 0 {
		val := mylog.Check2(ctx.transExpr(node.Inner[0]))

		// 指针无需显式初始化为空
	}
	if len(node.Inner) > 0 && node.Inner[len(node.Inner)-1].Kind == cast.FullComment {
		decl.Doc = ctx.transFullComment(node.Inner[len(node.Inner)-1], node.Name)
	}
	return decl, nil
}

func (ctx *TransCtx) transDecl(node *cast.Node) (decl ast.Decl, err error) {
	ctx.InNode(node)

	switch node.Kind {
	case cast.AccessSpecDecl:
		// ignore public, protected, private
	case cast.CXXConstructorDecl:
		decl = mylog.Check2(ctx.transCXXConstructorDecl(node))
	case cast.CXXDestructorDecl:
		decl = mylog.Check2(ctx.transCXXDestructorDecl(node))
	case cast.CXXMethodDecl:
		decl = mylog.Check2(ctx.transCXXMethodDecl(node))
	case cast.CXXRecordDecl:
		ctx.transCXXRecordDecl(node)
	case cast.EnumDecl:
		decl = mylog.Check2(ctx.transEnumDecl(node))
	case cast.FriendDecl:
		// ignore friend
	case cast.FunctionDecl:
		decl = mylog.Check2(ctx.transFunctionDecl(node))
	case cast.LinkageSpecDecl:
		ctx.transLinkageSpecDecl(node)
	case cast.NamespaceDecl:
		ctx.transNamespaceDecl(node)
	case cast.RecordDecl:
		decl = mylog.Check2(ctx.transRecordDecl(node))
	case cast.TypedefDecl:
		decl = mylog.Check2(ctx.transTypedefDecl(node))
	case cast.UsingDirectiveDecl:
		// ignore using namespace XXX;
	case cast.VarDecl:
		decl = mylog.Check2(ctx.transVarDecl(node))
	default:
		mylog.Check(ctx.newKindErr(node))
	}
	return
}
