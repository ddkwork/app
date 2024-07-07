package glfw

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
)

func TestGlfw(t *testing.T) {
	pkg := gengo.NewPackage("glfw",
		gengo.WithRemovePrefix(
			"glfw",
		))
	path := "D:\\fork\\glfw-master\\include\\GLFW\\glfw3.h"
	mylog.Check(pkg.Transform("glfw", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}
