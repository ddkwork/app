package glfw3_test

import (
	"testing"
	"unsafe"

	glfw3 "github.com/ddkwork/app/bindgen/project/glfw/tmp"
	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
)

func TestBind(t *testing.T) {
	windows.SetDllDirectory(".")
	glfw3.Init()
	mylog.Info("version", BytePointerToString(glfw3.GetVersionString()))
	defer glfw3.Terminate()
	w := glfw3.CreateWindow(200, 200, StringToBytePointer("hello word"), nil, nil)
	glfw3.MakeContextCurrent(w)
	for {
		glfw3.PollEvents()
		glfw3.SwapBuffers(w)
		if glfw3.WindowShouldClose(w) != 0 {
			glfw3.DestroyWindow(w)
			break
		}
	}
}

func StringToBytePointer(s string) *byte {
	bytes := []byte(s)
	ptr := &bytes[0]
	return ptr
}

func BytePointerToString(ptr *byte) string {
	var bytes []byte
	for *ptr != 0 {
		bytes = append(bytes, *ptr)
		ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
	}
	return string(bytes)
}
