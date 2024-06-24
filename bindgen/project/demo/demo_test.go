package libdemo

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"

	"github.com/ddkwork/golibrary/mylog"
)

func TestDemoDll(t *testing.T) {
	pkg := gengo.NewPackage("libdemo")
	path := "cpp\\library.h"
	mylog.Check(pkg.Transform("libdemo", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}
