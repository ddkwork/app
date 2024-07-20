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

const (
	COMMUNICATION_BUFFER_SIZE     = 256
	TCP_END_OF_BUFFER_CHARS_COUNT = 4
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

	pfn := func(msg *byte) {
		if msg == nil {
			println("msg is nil,callback not be called")
		}
		goData := (*[COMMUNICATION_BUFFER_SIZE + TCP_END_OF_BUFFER_CHARS_COUNT]byte)(unsafe.Pointer(msg))
		fmt.Println("Received data:", string(goData[:]))
		//assert.Equal(t, "TempMessage log callback buf test", string(goData[:]))
	}
	SetTextMessageCallback(unsafe.Pointer(reflect.ValueOf(pfn).Pointer()))
	ShowMessages(nil)
}
