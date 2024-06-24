package translate

import (
	cast "clang/ast"
	"go/ast"
	"strings"
	"utils"
)

func (ctx *TransCtx) transBlockCommandComment(node *cast.Node, name string) *ast.CommentGroup {
	if len(node.Inner) == 0 {
		utils.Err(ctx.newErr("block command comment enpty"))
		return &ast.CommentGroup{}
	}
	if node.Inner[0].Kind != cast.ParagraphComment {
		utils.Err(ctx.newKindErr(node.Inner[0]))
		return &ast.CommentGroup{}
	}
	pre := node.Name
	if pre == "brief" {
		pre = name + " :"
	}
	return ctx.transParagraphComment(node.Inner[0], pre)
}

func (ctx *TransCtx) transFullComment(node *cast.Node, name string) *ast.CommentGroup {
	cmnts := &ast.CommentGroup{}
	for _, n := range node.Inner {
		ctx.InNode(n)
		switch n.Kind {
		case cast.BlockCommandComment:
			list := ctx.transBlockCommandComment(n, name).List
			cmnts.List = append(cmnts.List, list...)
		case cast.ParagraphComment:
			list := ctx.transParagraphComment(n, name).List
			cmnts.List = append(cmnts.List, list...)
		case cast.ParamCommandComment:
			list := ctx.transParamCommandComment(n).List
			cmnts.List = append(cmnts.List, list...)
		default:
			utils.Err(ctx.newKindErr(n))
		}
	}
	return cmnts
}

func (ctx *TransCtx) transInlineCommandComment(node *cast.Node) string {
	return "\\" + node.Name + " " + strings.Join(node.Args, " ")
}

func (ctx *TransCtx) transParagraphComment(node *cast.Node, pre string) *ast.CommentGroup {
	if pre == "" {
		pre = "//"
	} else {
		pre = "// " + pre
	}
	cmnts := &ast.CommentGroup{}
	text := ""
	for _, n := range node.Inner {
		ctx.InNode(n)
		switch n.Kind {
		case cast.InlineCommandComment:
			text += ctx.transInlineCommandComment(n)
		case cast.TextComment:
			text += ctx.transTextComment(n)
		default:
			utils.Err(ctx.newKindErr(n))
		}

		if len(text) > 50 {
			cmnts.List = append(cmnts.List, &ast.Comment{Text: pre + text})
			pre = "//"
			text = ""
		}
	}
	if strings.TrimSpace(text) == "" {
		return cmnts
	}
	cmnts.List = append(cmnts.List, &ast.Comment{Text: pre + text})
	return cmnts
}

func (ctx *TransCtx) transParamCommandComment(node *cast.Node) *ast.CommentGroup {
	if len(node.Inner) == 0 {
		utils.Err(ctx.newErr("param command comment enpty"))
		return &ast.CommentGroup{}
	}
	if node.Inner[0].Kind != cast.ParagraphComment {
		utils.Err(ctx.newKindErr(node.Inner[0]))
		return &ast.CommentGroup{}
	}
	return ctx.transParagraphComment(node.Inner[0], node.Param)
}

func (ctx *TransCtx) transTextComment(node *cast.Node) string {
	return node.Text
}
