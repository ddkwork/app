package clang

import (
	"os"
	"path/filepath"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/golibrary/mylog"
)

type Options struct {
	ToolkitPath      string
	AdditionalParams []string
	Sources          []string
}

func (o *Options) ClangPath() string {
	if o.ToolkitPath != "" {
		if stat, e := os.Stat(o.ToolkitPath); e == nil && stat.IsDir() {
			return filepath.Join(o.ToolkitPath, "clang")
		} else {
			return o.ToolkitPath
		}
	}
	return "clang"
}

func (o *Options) ClangCommand(opt ...string) ([]byte, error) {
	c := make([]string, 0)
	c = append(c, o.ClangPath())
	c = append(c, opt...)
	c = append(c, o.AdditionalParams...)
	c = append(c, o.Sources...)
	return stream.RunCommandArgs(c...).Output.Bytes(), nil
}

func CreateAST(opt *Options) ([]byte, error) {
	return opt.ClangCommand(
		"-fsyntax-only",
		"-nobuiltininc",
		"-Xclang",
		"-ast-dump=json",
	)
}

func CreateLayoutMap(opt *Options) ([]byte, error) {
	return opt.ClangCommand(
		"-fsyntax-only",
		"-nobuiltininc",
		"-emit-llvm",
		"-Xclang",
		"-fdump-record-layouts",
		"-Xclang",
		"-fdump-record-layouts-complete",
	)
}

func Parse(opt *Options) (ast Node, layout *LayoutMap, err error) {
	// stream.RunCommand("clang -E -dM " + opt.Sources[0] + " > macros.log")
	res := mylog.Check2(CreateLayoutMap(opt))
	stream.WriteTruncate("astLayout.log", res)
	layout = mylog.Check2(ParseLayoutMap(res))

	res = mylog.Check2(CreateAST(opt))
	stream.WriteTruncate("ast.json", res)
	ast = mylog.Check2(ParseAST(res))

	return ast, layout, nil
}
