package callback

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"

	"github.com/ddkwork/golibrary/mylog"
)

func TestDemoDll(t *testing.T) {
	pkg := gengo.NewPackage("callback")
	path := "src/callback.h"
	mylog.Check(pkg.Transform("callback", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("."))
	//Hello()

	pfn := func(msg *byte) {
		if msg == nil {
			println("msg is nil,callback not be called")
		}
		goData := (*[260]byte)(unsafe.Pointer(&msg))
		fmt.Println("Received data:", string(goData[:]))
	}
	SetTextMessageCallback(unsafe.Pointer(reflect.ValueOf(pfn).Pointer()))
	ShowMessages(nil)
}
