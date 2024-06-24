package translate

import (
	cast "clang/ast"
	"fmt"
	"go/ast"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

var voidType = &ast.SelectorExpr{X: &ast.Ident{Name: "unsafe"}, Sel: &ast.Ident{Name: "ArbitraryType"}}

var builtinType = map[string]ast.Expr{
	"void":               voidType,
	"char":               &ast.Ident{Name: "byte"},
	"signed char":        &ast.Ident{Name: "int8"},
	"unsigned char":      &ast.Ident{Name: "uint8"},
	"short":              &ast.Ident{Name: "int16"},
	"unsigned short":     &ast.Ident{Name: "uint16"},
	"int":                &ast.Ident{Name: "int"},
	"unsigned int":       &ast.Ident{Name: "uint"},
	"long":               &ast.Ident{Name: "int32"},
	"unsigned long":      &ast.Ident{Name: "uint32"},
	"long long":          &ast.Ident{Name: "int64"},
	"unsigned long long": &ast.Ident{Name: "uint64"},
	"float":              &ast.Ident{Name: "float32"},
	"double":             &ast.Ident{Name: "float64"},
	"long double":        &ast.Ident{Name: "float64"},
	"_Complex float":     &ast.Ident{Name: "complex64"},
	"_Complex double":    &ast.Ident{Name: "complex128"},
}

func (ctx *TransCtx) transBuiltinType(node *cast.Node) (t ast.Expr, err error) {
	if node.Type == nil {
		mylog.Check(ctx.newErr("unknown builtin type"))
		return
	}
	ok := false
	if t, ok = builtinType[node.Type.QualType]; !ok {
		mylog.Check(ctx.newErr("unknown builtin type: ", node.Type.QualType))
	}
	return
}

func (ctx *TransCtx) transElaboratedType(node *cast.Node) (ast.Expr, error) {
	if node.OwnedTagDecl != nil {
		if node.OwnedTagDecl.Name != "" {
			t := &ast.Ident{Name: node.OwnedTagDecl.Name}
			return t, nil
		}
		if node.OwnedTagDecl.Kind == cast.EnumDecl {
			return &ast.Ident{Name: "int"}, nil
		}
		id := node.ID.ToInt()
		if id > 0 && ctx.nodes[id].data != nil {
			if t, ok := ctx.nodes[id].data.(ast.Expr); ok {
				return t, nil
			}
			utils.Log(utils.LL_Verbose, "node data unexpect: %+v", ctx.nodes[id].data)
		}
	}
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("unknown base type")
	}
	return ctx.transType(node.Inner[0])
}

func (ctx *TransCtx) transFunctionProtoType(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("function without return type")
	}
	t := &ast.FuncType{Params: &ast.FieldList{}}

	t0 := mylog.Check2(ctx.transType(node.Inner[0]))

	if t0 != nil {
		field := &ast.Field{Type: t0}
		t.Results = &ast.FieldList{List: []*ast.Field{field}}
	}

	if len(node.Inner) == 1 {
		return t, nil
	}
	for _, parm := range node.Inner[1:] {
		t0 = mylog.Check2(ctx.transType(parm))

		field := &ast.Field{Type: t0}
		if parm.Name != "" {
			field.Names = []*ast.Ident{{Name: parm.Name}}
		}
		t.Params.List = append(t.Params.List, field)
	}
	return t, nil
}

func (ctx *TransCtx) transParenType(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("paren without base type")
	}
	return ctx.transType(node.Inner[0])
}

func (ctx *TransCtx) transPointerType(node *cast.Node) (t ast.Expr, err error) {
	if len(node.Inner) == 0 {
		mylog.Check(ctx.newErr("pointer without base type"))
		return
	}
	t = mylog.Check2(ctx.transType(node.Inner[0]))

	// void* -> unsafe.Pointer
	if t == voidType {
		t = &ast.SelectorExpr{X: &ast.Ident{Name: "unsafe"}, Sel: &ast.Ident{Name: "Pointer"}}
		return
	}
	// void(*)() -> func()
	if _, ok := t.(*ast.FuncType); ok && node.Inner[0].Kind == cast.ParenType {
		return
	}
	t = &ast.StarExpr{X: t}
	return
}

func (ctx *TransCtx) transQualType(node *cast.Node) (t ast.Expr, err error) {
	if len(node.Inner) == 0 {
		mylog.Check(ctx.newErr("qualifiers without base type"))
		return
	}
	utils.Log(utils.LL_Debug, "%s.qualifiers=%s", cast.QualType, node.Qualifiers)
	t = mylog.Check2(ctx.transType(node.Inner[0]))
	return
}

func (ctx *TransCtx) transRecordType(node *cast.Node) (ast.Expr, error) {
	if node.Decl == nil {
		return nil, ctx.newErr("record without base type")
	}
	id := node.Decl.ID.ToInt()
	if id > 0 && ctx.nodes[id].data != nil {
		if t, ok := ctx.nodes[id].data.(ast.Expr); ok {
			return t, nil
		}
		return nil, ctx.newErr(fmt.Sprintf("node data unexpect: %+v", ctx.nodes[id].data))
	}
	if node.Type != nil {
		t := mylog.Check2(NewTypeParse(node.Type.QualType).Get())

		return t, nil
	}
	return nil, ctx.newErr("unknow type: ", node.Name)
}

func (ctx *TransCtx) transTempSpecialType(node *cast.Node) (ast.Expr, error) {
	if len(node.Inner) == 0 {
		return nil, ctx.newErr("template without argument")
	}
	var list []ast.Expr
	for _, n := range node.Inner {
		ctx.InNode(n)
		switch n.Kind {
		case cast.RecordType:
		// ignore template define
		case cast.TemplateArgument:
			if len(n.Inner) == 0 {
				utils.Err(ctx.newErr("template argument without base type"))
				break
			}
			t := mylog.Check2(ctx.transType(n.Inner[0]))

		default:
			utils.Err(ctx.newKindErr(n))
		}
	}

	t := &ast.Ident{Name: node.TemplateName}
	if len(list) == 1 {
		return &ast.IndexExpr{X: t, Index: list[0]}, nil
	} else if len(list) > 1 {
		return &ast.IndexListExpr{X: t, Indices: list}, nil
	}
	return t, nil
}

func (ctx *TransCtx) transTypedefType(node *cast.Node) (ast.Expr, error) {
	if node.Decl != nil {
		if t, ok := ctx.nodes[node.Decl.ID.ToInt()].data.(ast.Expr); ok {
			return t, nil
		}
		if node.Decl.Name != "" {
			var t ast.Expr = &ast.Ident{Name: node.Decl.Name}
			ctx.setNodeData(node.Decl.ID.ToInt(), t)
			return t, nil
		}
	}
	if node.Type != nil {
		return &ast.Ident{Name: node.Type.QualType}, nil
	}
	if len(node.Inner) > 0 {
		return ctx.transType(node.Inner[0])
	}
	return nil, ctx.newErr("unknown typedef type")
}

func (ctx *TransCtx) transType(node *cast.Node) (t ast.Expr, err error) {
	ctx.InNode(node)
	id := node.ID.ToInt()
	if id > 0 && ctx.nodes[id].data != nil {
		ok := false
		if t, ok = ctx.nodes[id].data.(ast.Expr); ok {
			return
		}
		utils.Log(utils.LL_Error, "node data unexpect: %+v", ctx.nodes[id].data)
	}

	switch node.Kind {
	case cast.BuiltinType:
		t = mylog.Check2(ctx.transBuiltinType(node))
	case cast.ElaboratedType:
		t = mylog.Check2(ctx.transElaboratedType(node))
	case cast.FunctionProtoType:
		t = mylog.Check2(ctx.transFunctionProtoType(node))
	case cast.ParenType:
		t = mylog.Check2(ctx.transParenType(node))
	case cast.PointerType:
		t = mylog.Check2(ctx.transPointerType(node))
	case cast.QualType:
		t = mylog.Check2(ctx.transQualType(node))
	case cast.RecordType:
		t = mylog.Check2(ctx.transRecordType(node))
	case cast.TempSpecialType:
		t = mylog.Check2(ctx.transTempSpecialType(node))
	case cast.TypedefType:
		t = mylog.Check2(ctx.transTypedefType(node))
	default:
		mylog.Check(ctx.newKindErr(node))
		return
	}
	if err == nil {
		ctx.setNodeData(id, t)
	}
	return
}
