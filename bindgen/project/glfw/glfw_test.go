package glfw

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
)

func TestGlfw(t *testing.T) {
	t.Skip()
	pkg := gengo.NewPackage("glfw3",
		gengo.WithRemovePrefix(
			"glfw",
			"gl",
		))
	path := "include\\GLFW\\glfw3.h"
	mylog.Check(pkg.Transform("glfw3", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}
