package parser

import (
	"bytes"
	"clang/ast"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

type Config struct {
	Json       string
	Flags      []string
	Stderr     bool
	SkipOthers bool
}

func DumpAST(filename string, conf *Config) (result []byte, warning []byte, err error) {
	if conf == nil {
		conf = &Config{}
	}
	args := []string{"-Xclang", "-ast-dump=json", "-fsyntax-only"}
	if len(conf.Flags) != 0 {
		args = append(args, conf.Flags...)
		if strings.Contains(strings.Join(conf.Flags, " "), "-nostdinc") {
			incDir := filepath.Join(filepath.Dir(os.Args[0]), "../include")
			args = append(args, "-I"+incDir)
		}
	}
	args = append(args, filename)

	utils.Log(utils.LL_Info, "==> run CMD: clang %s", args)
	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd := exec.Command("clang", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	mylog.Check(cmd.Run())
	return outBuf.Bytes(), errBuf.Bytes(), err
}

func ParseFile(filename string, conf *Config) (*ast.Node, error) {
	if conf == nil {
		conf = &Config{}
	}
	out, warning := mylog.Check3(DumpAST(filename, conf))
	if len(out) == 0 {
		if len(warning) > 0 {
			mylog.Check(errors.New(string(warning)))
		}
		return nil, err
	}
	if conf.Stderr && len(warning) > 0 {
		_, _ = os.Stderr.Write(warning)
	}
	if conf.Json != "" && !conf.SkipOthers {
		mylog.Check(os.WriteFile(conf.Json, out, 644))
	}

	utils.Log(utils.LL_Info, "==> unmarshal json")
	var file ast.Node
	mylog.Check(json.Unmarshal(out, &file))

	if conf.SkipOthers {
		utils.Log(utils.LL_Info, "==> skip other files")
		skipOthers(filename, &file)
		if conf.Json != "" {
			out = mylog.Check2(json.Marshal(&file))
		}
	}
	return &file, nil
}

func LoadJsonFile(filename string) (*ast.Node, error) {
	utils.Log(utils.LL_Info, "==> open json file")
	out := mylog.Check2(os.ReadFile(filename))

	utils.Log(utils.LL_Info, "==> unmarshal json")
	var file ast.Node
	mylog.Check(json.Unmarshal(out, &file))

	return &file, nil
}

func skipOthers(filename string, file *ast.Node) {
	if file.Kind != ast.TranslationUnitDecl {
		utils.Log(utils.LL_Error, "unknown node kind: %s", file.Kind)
		return
	}

	name := filepath.Base(filename)
	if index := strings.LastIndex(name, "."); index > 0 {
		name = name[:index+1]
	}

	var subs []*ast.Node
	path := ""
	for _, sub := range file.Inner {
		if sub.Loc != nil && sub.Loc.File != "" {
			path = sub.Loc.File
		} else if sub.Range != nil && sub.Range.Begin.Loc != nil && sub.Range.Begin.Loc.File != "" {
			path = sub.Range.Begin.Loc.File
		}

		if strings.HasPrefix(filepath.Base(path), name) {
			subs = append(subs, sub)
		}
	}

	file.Inner = subs
}
