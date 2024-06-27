package libdemo

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"

	"github.com/ddkwork/golibrary/mylog"
)

func TestSkia(t *testing.T) {
	t.Skip()
	pkg := gengo.NewPackage("skia")
	path := "sk_capi.h"
	mylog.Check(pkg.Transform("skia", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}
