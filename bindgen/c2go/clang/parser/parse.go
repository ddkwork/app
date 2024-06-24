package parser

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/ddkwork/app/bindgen/c2go/clang/ast"
	"github.com/ddkwork/golibrary/mylog"
	jsoniter "github.com/json-iterator/go"
)

type Mode uint

// -----------------------------------------------------------------------------

type ParseError struct {
	Err    error
	Stderr []byte
}

func (p *ParseError) Error() string {
	if len(p.Stderr) > 0 {
		return string(p.Stderr)
	}
	return p.Err.Error()
}

// -----------------------------------------------------------------------------

type Config struct {
	Json   *[]byte
	Flags  []string
	Stderr bool
}

func DumpAST(filename string, conf *Config) (result []byte, warning []byte, err error) {
	if conf == nil {
		conf = new(Config)
	}
	skiperr := strings.HasSuffix(filename, "vfprintf.c.i")
	stdout := NewPagedWriter()
	stderr := new(bytes.Buffer)
	args := []string{"-Xclang", "-ast-dump=json", "-fsyntax-only", filename}
	if len(conf.Flags) != 0 {
		args = append(conf.Flags, args...)
	}
	cmd := exec.Command("clang", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	if conf.Stderr && !skiperr {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = stderr
	}
	mylog.Check(cmd.Run())
	errmsg := stderr.Bytes()
	if err != nil && !skiperr {
		return nil, nil, &ParseError{Err: err, Stderr: errmsg}
	}
	return stdout.Bytes(), errmsg, nil
}

// -----------------------------------------------------------------------------

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ParseFileEx(filename string, mode Mode, conf *Config) (file *ast.Node, warning []byte, err error) {
	out, warning := mylog.Check3(DumpAST(filename, conf))

	if conf != nil && conf.Json != nil {
		*conf.Json = out
	}
	file = new(ast.Node)
	mylog.Check(json.Unmarshal(out, file))

	return
}

func ParseFile(filename string, mode Mode) (file *ast.Node, warning []byte, err error) {
	return ParseFileEx(filename, mode, nil)
}

// -----------------------------------------------------------------------------
