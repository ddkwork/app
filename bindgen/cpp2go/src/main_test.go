package main

import (
	"go/parser"
	"go/token"
	"os"
	"testing"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

func TestParseGoFile(t *testing.T) {
	fset := token.NewFileSet()
	file := mylog.Check2(parser.ParseFile(fset, "test/main.go", nil, parser.ParseComments))

	// utils.Log(utils.LL_Info, "%+v", file)
	buf := mylog.Check2(utils.Dump("", file))
	mylog.Check(os.WriteFile("test/dump.txt", buf, 644))
	utils.Err(err)
}
