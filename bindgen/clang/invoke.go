package clang

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"

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
	cmd := exec.Command(o.ClangPath(), opt...)
	cmd.Args = append(cmd.Args, o.AdditionalParams...)
	cmd.Args = append(cmd.Args, o.Sources...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	mylog.Check(cmd.Run())

	return buf.Bytes(), nil

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
	errg := &errgroup.Group{}
	errg.Go(func() error {
		res, e := CreateAST(opt)
		if e != nil {
			return e
		}
		ast, e = ParseAST(res)
		return e
	})
	errg.Go(func() error {
		res, e := CreateLayoutMap(opt)
		if e != nil {
			return e
		}
		layout, e = ParseLayoutMap(res)
		return e
	})
	if mylog.Check(errg.Wait()); err != nil {
		return nil, nil, err
	}
	return ast, layout, nil
}

func Parse_(opt *Options) (ast Node, layout *LayoutMap, err error) {
	// stream.RunCommand("clang -E -dM " + opt.Sources[0] + " > macros.log")
	res := mylog.Check2(CreateLayoutMap(opt))
	stream.WriteTruncate("astLayout.log", res)
	layout = mylog.Check2(ParseLayoutMap(res))

	res = mylog.Check2(CreateAST(opt))
	stream.WriteTruncate("ast.json", res)
	ast = mylog.Check2(ParseAST(res))

	return ast, layout, nil
}
