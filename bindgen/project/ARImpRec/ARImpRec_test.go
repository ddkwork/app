package ARImpRec

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"

	"github.com/ddkwork/golibrary/mylog"
)

func TestARImpRec(t *testing.T) {
	t.Skip()
	pkg := gengo.NewPackage("ARImpRec")
	mylog.Check(pkg.Transform("ARImpRec", &clang.Options{
		Sources:          []string{"ARImpRec.h"},
		AdditionalParams: []string{},
	}))
	mylog.Check(pkg.WriteToDir("./tmp"))
}
