package glfw

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
)

func TestGlfw(t *testing.T) {
	pkg := gengo.NewPackage("libdemo")
	path := "D:\\fork\\glfw-master\\include\\GLFW\\glfw3.h"
	mylog.Check(pkg.Transform("libdemo", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}
