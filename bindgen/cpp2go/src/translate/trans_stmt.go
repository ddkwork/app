package translate

import (
	cast "clang/ast"
	"go/ast"
	"go/token"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	breakStmt       = &ast.BranchStmt{Tok: token.BREAK}
	continueStmt    = &ast.BranchStmt{Tok: token.CONTINUE}
	fallthroughStmt = &ast.BranchStmt{Tok: token.FALLTHROUGH}
	nullStmt        = &ast.EmptyStmt{}
)

// newDefineStmt : 新建 name := val 语句
func newDefineStmt(name, val ast.Expr) ast.Stmt {
	return &ast.AssignStmt{Lhs: []ast.Expr{name}, Tok: token.DEFINE, Rhs: []ast.Expr{val}}
}

func (ctx *TransCtx) transAssignStmt(node *cast.Node) (ast.Stmt, error) {
	if node.OpCode != cast.Assign {
		return nil, ctx.newErr("unsupport op code: ", node.OpCode)
	}
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown left or right value for binary operator")
	}
	left := mylog.Check2(ctx.transExpr(node.Inner[0]))

	right := mylog.Check2(ctx.transExpr(node.Inner[1]))

	stmt := &ast.AssignStmt{Lhs: []ast.Expr{left}, Tok: token.ASSIGN, Rhs: []ast.Expr{right}}
	return stmt, nil
}

// transCXXCatchStmt : C++ catch语句
func (ctx *TransCtx) transCXXCatchStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown catch statement")
	}

	stmt := &ast.CaseClause{}
	for _, n := range node.Inner {
		ctx.InNode(n)
		switch n.Kind {
		case cast.VarDecl:
			decl := mylog.Check2(ctx.transVarDecl(n))

			// name := err

		case cast.CompoundStmt:
			list := ctx.transCompoundStmt(n).List
			stmt.Body = append(stmt.Body, list...)
		default:
			utils.Err(ctx.newKindErr(node))
		}
	}
	return stmt, nil
}

// transCXXTryStmt : C++ try语句
func (ctx *TransCtx) transCXXTryStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown try statement")
	}
	var catchs []ast.Stmt
	var normal ast.Stmt
	for _, n := range node.Inner {
		ctx.InNode(n)
		if n.Kind == cast.CXXCatchStmt {
			catch := mylog.Check2(ctx.transCXXCatchStmt(n))
		} else {
			stmt := mylog.Check2(ctx.transStmt(n))
		}
	}

	r := &ast.Ident{Name: "r"}
	// r := recover()
	rcvDef := newDefineStmt(r, &ast.CallExpr{Fun: &ast.Ident{Name: "recover"}})
	// if r == nil { return }
	rcvNil := &ast.IfStmt{Cond: &ast.BinaryExpr{X: r, Op: token.EQL, Y: nullExpr}, Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}}}
	defers := []ast.Stmt{rcvDef, rcvNil}
	if len(catchs) > 0 {
		// switch err := r.(type) { ... }
		errDef := newDefineStmt(&ast.Ident{Name: "err"}, &ast.TypeAssertExpr{X: r})
		rcvSw := &ast.TypeSwitchStmt{Assign: errDef, Body: &ast.BlockStmt{List: catchs}}
		defers = append(defers, rcvSw)
	}

	// defer func() { ... }()
	fn := &ast.FuncLit{Type: &ast.FuncType{Params: &ast.FieldList{}}, Body: &ast.BlockStmt{List: defers}}
	deferStmt := &ast.DeferStmt{Call: &ast.CallExpr{Fun: fn}}
	if normal == nil {
		return deferStmt, nil
	}

	stmt := &ast.BlockStmt{List: []ast.Stmt{deferStmt}}
	if blk, ok := normal.(*ast.BlockStmt); ok {
		stmt.List = append(stmt.List, blk.List...)
	} else {
		stmt.List = append(stmt.List, normal)
	}
	return stmt, nil
}

func (ctx *TransCtx) transCaseStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown judging condition")
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[0]))

	stmt := &ast.CaseClause{List: []ast.Expr{cond}}
	for _, n := range node.Inner[1:] {
		st := mylog.Check2(ctx.transStmt(n))
	}

	if len(stmt.Body) == 1 {
		if st, ok := stmt.Body[0].(*ast.CaseClause); ok {
			if len(st.List) > 0 {
				stmt.List = append(stmt.List, st.List...)
				stmt.Body = st.Body
			} else {
				stmt.Body = []ast.Stmt{fallthroughStmt, st}
			}
		}
	} else if len(stmt.Body) > 1 {
		if st, ok := stmt.Body[len(stmt.Body)-1].(*ast.CaseClause); ok {
			stmt.Body = append(stmt.Body[:len(stmt.Body)-1], fallthroughStmt, st)
		}
	}
	return stmt, nil
}

func (ctx *TransCtx) transCompoundAssignStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown left or right value for binary operator")
	}
	left := mylog.Check2(ctx.transExpr(node.Inner[0]))

	right := mylog.Check2(ctx.transExpr(node.Inner[1]))

	op, ok := operatorCodeAssignMap[node.OpCode]
	if !ok {
		return nil, ctx.newErr("unsupport op code: ", node.OpCode)
	}
	stmt := &ast.AssignStmt{Lhs: []ast.Expr{left}, Tok: op, Rhs: []ast.Expr{right}}
	return stmt, nil
}

// transCompoundStmt : 大括号 { }
func (ctx *TransCtx) transCompoundStmt(node *cast.Node) *ast.BlockStmt {
	stmts := &ast.BlockStmt{}
	if len(node.Inner) == 0 {
		return stmts
	}
	for _, n := range node.Inner {
		stmt := mylog.Check2(ctx.transStmt(n))
	}
	return stmts
}

func (ctx *TransCtx) transDeclStmt(node *cast.Node) (ast.Stmt, error) {
	var decls []ast.Decl
	for _, n := range node.Inner {
		decl := mylog.Check2(ctx.transDecl(n))

		decls = append(decls, decl)
	}
	if len(decls) == 0 {
		return nil, ctx.newErr("empty declare")
	}

	decl := decls[0]
	if len(decls) == 1 {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.VAR && len(gen.Specs) == 1 {
			// (var a = uint(0)) -> (a := uint(0))
			if spec, ok := gen.Specs[0].(*ast.ValueSpec); ok && spec.Type == nil && len(spec.Names) == 1 {
				stmt := &ast.AssignStmt{Lhs: []ast.Expr{spec.Names[0]}, Tok: token.DEFINE, Rhs: spec.Values}
				return stmt, nil
			}
		}
	} else {
		var specs []ast.Spec
		var names []*ast.Ident
		var values []ast.Expr
		var t ast.Expr
		for _, decl = range decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				utils.Log(utils.LL_Error, "unknow decl type: %+v", decl)
				continue
			}
			specs = append(specs, gen.Specs...)

			for _, spec := range gen.Specs {
				valSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					break
				}
				names = append(names, valSpec.Names...)
				if valSpec.Type == nil {
					values = append(values, valSpec.Values...)
				} else {
					t = valSpec.Type
				}
			}
		}

		if len(specs) > 0 {
			if len(specs) == len(names) {
				if len(names) == len(values) {
					// (var (a = uint(0); b = 1)) -> (a, b := uint(0), 1)
					stmt := &ast.AssignStmt{Tok: token.DEFINE, Rhs: values}
					for _, name := range names {
						stmt.Lhs = append(stmt.Lhs, name)
					}
					return stmt, nil
				} else if len(values) > 0 {
					decl = &ast.GenDecl{Tok: token.VAR, Specs: specs}
				} else {
					// (var (a int; b int)) -> (var a, b int)
					spec := &ast.ValueSpec{Names: names, Type: t}
					decl = &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{spec}}
				}
			} else {
				decl = &ast.GenDecl{Tok: token.VAR, Specs: specs}
			}
		}
	}

	stmt := &ast.DeclStmt{Decl: decl}
	return stmt, nil
}

func (ctx *TransCtx) transDefaultStmt(node *cast.Node) (ast.Stmt, error) {
	stmt := &ast.CaseClause{}
	for _, n := range node.Inner {
		st := mylog.Check2(ctx.transStmt(n))
	}
	if len(stmt.Body) > 0 {
		if st, ok := stmt.Body[len(stmt.Body)-1].(*ast.CaseClause); ok {
			stmt.Body = append(stmt.Body[:len(stmt.Body)-1], fallthroughStmt, st)
		}
	}
	return stmt, nil
}

func (ctx *TransCtx) transDoStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 执行语句
	stat := mylog.Check2(ctx.transStmt(node.Inner[0]))

	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[1]))

	if c, ok := cond.(*ast.BasicLit); ok && c.Kind == token.INT && c.Value == "0" {
		// do {...} while (0) -> {...}
		return stat, nil
	}

	flag := &ast.Ident{Name: "__flag"}
	stmt := &ast.ForStmt{Cond: flag}
	stmt.Init = &ast.AssignStmt{Lhs: []ast.Expr{flag}, Tok: token.ASSIGN, Rhs: []ast.Expr{&ast.Ident{Name: "true"}}}
	stmt.Post = &ast.AssignStmt{Lhs: []ast.Expr{flag}, Tok: token.ASSIGN, Rhs: []ast.Expr{cond}}
	if body, ok := stat.(*ast.BlockStmt); ok {
		stmt.Body = body
	} else {
		stmt.Body = &ast.BlockStmt{List: []ast.Stmt{stat}}
	}
	return stmt, nil
}

func (ctx *TransCtx) transForStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 5 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 初始化语句
	init := mylog.Check2(ctx.transStmt(node.Inner[0]))
	if err != nil && node.Inner[0].Kind != "" {
		utils.Err(err)
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[2]))
	if err != nil && node.Inner[2].Kind != "" {
		utils.Err(err)
		cond = &ast.Ident{Name: "false"}
	}
	// 步进语句
	post := mylog.Check2(ctx.transStmt(node.Inner[3]))
	if err != nil && node.Inner[3].Kind != "" {
		utils.Err(err)
	}
	// 执行语句
	stat := mylog.Check2(ctx.transStmt(node.Inner[4]))

	stmt := &ast.ForStmt{Init: init, Cond: cond, Post: post}
	if body, ok := stat.(*ast.BlockStmt); ok {
		stmt.Body = body
	} else {
		stmt.Body = &ast.BlockStmt{List: []ast.Stmt{stat}}
	}
	return stmt, nil
}

func (ctx *TransCtx) transGotoStmt(node *cast.Node) (ast.Stmt, error) {
	label, ok := ctx.nodes[node.TargetLabelDeclId.ToInt()].data.(*ast.LabeledStmt)
	if !ok || label == nil {
		stmt := &ast.BranchStmt{Tok: token.GOTO, Label: &ast.Ident{Name: "Unknow_" + string(node.TargetLabelDeclId)}}
		return stmt, nil
	}
	stmt := &ast.BranchStmt{Tok: token.GOTO, Label: label.Label}
	return stmt, nil
}

func (ctx *TransCtx) transIfStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[0]))

	// 成功分支
	stat := mylog.Check2(ctx.transStmt(node.Inner[1]))

	body, ok := stat.(*ast.BlockStmt)
	if !ok {
		body = &ast.BlockStmt{List: []ast.Stmt{stat}}
	}
	stmt := &ast.IfStmt{Cond: cond, Body: body}
	if len(node.Inner) > 2 {
		// 失败分支
		stat = mylog.Check2(ctx.transStmt(node.Inner[2]))
	}
	return stmt, nil
}

func (ctx *TransCtx) transLabelStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown label statement")
	}
	stmt := &ast.LabeledStmt{Label: &ast.Ident{Name: node.Name}}
	ctx.setNodeData(node.DeclId.ToInt(), stmt)

	stmt.Stmt = mylog.Check2(ctx.transStmt(node.Inner[0]))

	return stmt, nil
}

func (ctx *TransCtx) transReturnStmt(node *cast.Node) (ast.Stmt, error) {
	stmt := &ast.ReturnStmt{}
	if len(node.Inner) > 0 {
		expr := mylog.Check2(ctx.transExpr(node.Inner[0]))

		stmt.Results = append(stmt.Results, expr)
	}
	return stmt, nil
}

func (ctx *TransCtx) transSwitchStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[0]))

	// 条件分支
	stat := mylog.Check2(ctx.transStmt(node.Inner[1]))

	body, ok := stat.(*ast.BlockStmt)
	if !ok {
		body = &ast.BlockStmt{List: []ast.Stmt{stat}}
	}
	stmt := &ast.SwitchStmt{Tag: cond, Body: body}
	return stmt, nil
}

func (ctx *TransCtx) transWhileStmt(node *cast.Node) (ast.Stmt, error) {
	if len(node.Inner) < 2 {
		return nil, ctx.newErr("unknown judging condition or statement")
	}
	// 判断条件
	cond := mylog.Check2(ctx.transExpr(node.Inner[0]))

	// 执行语句
	stat := mylog.Check2(ctx.transStmt(node.Inner[1]))

	stmt := &ast.ForStmt{Cond: cond}
	if body, ok := stat.(*ast.BlockStmt); ok {
		stmt.Body = body
	} else {
		stmt.Body = &ast.BlockStmt{List: []ast.Stmt{stat}}
	}
	return stmt, nil
}

func (ctx *TransCtx) transStmt(node *cast.Node) (ast.Stmt, error) {
	ctx.InNode(node)

	switch node.Kind {
	case cast.BinaryOperator:
		return ctx.transAssignStmt(node)
	case cast.BreakStmt:
		return breakStmt, nil
	case cast.CXXTryStmt:
		return ctx.transCXXTryStmt(node)
	case cast.CaseStmt:
		return ctx.transCaseStmt(node)
	case cast.CompoundAssignOperator:
		return ctx.transCompoundAssignStmt(node)
	case cast.CompoundStmt:
		return ctx.transCompoundStmt(node), nil
	case cast.ContinueStmt:
		return continueStmt, nil
	case cast.DeclStmt:
		return ctx.transDeclStmt(node)
	case cast.DefaultStmt:
		return ctx.transDefaultStmt(node)
	case cast.DoStmt:
		return ctx.transDoStmt(node)
	case cast.ForStmt:
		return ctx.transForStmt(node)
	case cast.GotoStmt:
		return ctx.transGotoStmt(node)
	case cast.IfStmt:
		return ctx.transIfStmt(node)
	case cast.LabelStmt:
		return ctx.transLabelStmt(node)
	case cast.NullStmt:
		return nullStmt, nil
	case cast.ReturnStmt:
		return ctx.transReturnStmt(node)
	case cast.SwitchStmt:
		return ctx.transSwitchStmt(node)
	case cast.WhileStmt:
		return ctx.transWhileStmt(node)
	default:
		return ctx.transExprStmt(node)
	}
}

func (ctx *TransCtx) transExprStmt(node *cast.Node) (ast.Stmt, error) {
	expr := mylog.Check2(ctx.transExpr(node))

	// ++/-- 转换
	if unary, ok := expr.(*ast.UnaryExpr); ok {
		if unary.Op == token.INC || unary.Op == token.DEC {
			stmt := &ast.IncDecStmt{X: unary.X, Tok: unary.Op}
			return stmt, nil
		}
	}

	stmt := &ast.ExprStmt{X: expr}
	return stmt, nil
}
