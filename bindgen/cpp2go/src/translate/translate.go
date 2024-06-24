package translate

import (
	cast "clang/ast"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

// Keep these in sync with go/format/format.go.
const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent | printerNormalizeNumbers

	// printerNormalizeNumbers means to canonicalize number literal prefixes
	// and exponents while printing. See https://golang.org/doc/go1.13#gofmt.
	//
	// This value is defined in go/printer specifically for go/format and cmd/gofmt.
	printerNormalizeNumbers = 1 << 30
)

var goFormatConfig = printer.Config{Mode: printerMode, Tabwidth: tabWidth}

type nodeObj struct {
	node *cast.Node
	data any
}

type TransCtx struct {
	cFile   *cast.Node
	cfg     Config
	nodes   map[uint64]nodeObj
	fileSet *token.FileSet
	goFile  *ast.File

	loc cast.Loc
}

type Config struct {
	Package    string
	FileName   string // C/C++文件名
	SkipOthers bool
}

func NewTranslate(cFile *cast.Node, cfg *Config) *TransCtx {
	if cfg == nil {
		cfg = &Config{Package: "main"}
	}
	return &TransCtx{cFile: cFile, cfg: *cfg, nodes: make(map[uint64]nodeObj), fileSet: token.NewFileSet()}
}

func (ctx *TransCtx) InNode(node *cast.Node) {
	if node.ID != "" {
		id := node.ID.ToInt()
		if _, ok := ctx.nodes[id]; !ok {
			ctx.nodes[id] = nodeObj{node: node}
		}
	}

	loc := node.Loc
	if loc == nil && node.Range != nil {
		loc = node.Range.Begin.Loc
	}
	if loc != nil {
		if loc.File != "" {
			ctx.loc.File = loc.File
			ctx.loc.IncludedFrom = loc.IncludedFrom
		}
		if loc.Line > 0 {
			ctx.loc.Line = loc.Line
		}
		if loc.Col > 0 {
			ctx.loc.Col = loc.Col
		}
	}
}

func (ctx *TransCtx) Run() error {
	ctx.goFile = &ast.File{Package: 1}
	ctx.goFile.Name = &ast.Ident{NamePos: 9, Name: ctx.cfg.Package}

	ctx.InNode(ctx.cFile)
	if ctx.cFile.Kind != cast.TranslationUnitDecl {
		return ctx.newErr("unexpect: ", cast.TranslationUnitDecl)
	}

	name := filepath.Base(ctx.cfg.FileName)
	if index := strings.LastIndex(name, "."); index > 0 {
		name = name[:index+1]
	}

	for _, node := range ctx.cFile.Inner {
		// 忽略编译器自带的内建声明
		if node.IsImplicit {
			ctx.InNode(node)
			utils.Log(utils.LL_Debug, "ignore builtin AST: %s", node.Name)
			continue
		}

		// 跳过通过 include 引入的其他文件
		if ctx.cfg.SkipOthers {
			ctx.InNode(node)
			file := ctx.loc.File
			if file == "" || !strings.HasPrefix(filepath.Base(file), name) {
				utils.Log(utils.LL_Debug, "ignore include file AST: %s", node.Name)
				continue
			}
		}

		decl := mylog.Check2(ctx.transDecl(node))

	}
	return nil
}

func (ctx *TransCtx) Save(file string) error {
	f := mylog.Check2(os.Create(file))

	defer f.Close()

	ast.SortImports(ctx.fileSet, ctx.goFile)
	return goFormatConfig.Fprint(f, ctx.fileSet, ctx.goFile)
}

func (ctx *TransCtx) newErr(v ...any) error {
	return newTransErr(&ctx.loc, v...)
}

func (ctx *TransCtx) newKindErr(node *cast.Node) error {
	return newTransErr(&ctx.loc, "unsupport node kind: ", node.Kind)
}

func (ctx *TransCtx) setNodeData(id uint64, data any) {
	if id == 0 {
		return
	}
	ctx.nodes[id] = nodeObj{node: ctx.nodes[id].node, data: data}
}
