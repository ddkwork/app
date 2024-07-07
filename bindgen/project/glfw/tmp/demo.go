package main

import (
	_ "embed"
	"runtime"
	"unsafe"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/golibrary/mylog"
)

//go:embed glfw3.dll
var dll []byte

func init() {
	runtime.LockOSThread()
	path := "C:\\Windows\\glfw3.dll"
	if !stream.IsFilePath(path) {
		stream.WriteTruncate(path, dll)
	}
	mylog.Check2(GengoLibrary.LoadEmbed(dll))
}

func main() {
	// path := "D:\\workspace\\workspace\\app\\bindgen\\project\\glfw\\tmp"
	// mylog.Check(windows.SetDllDirectory(path))
	Init()
	mylog.Info("version", BytePointerToString(GetVersionString()))
	defer Terminate()
	w := CreateWindow(200, 200, StringToBytePointer("hello word"), nil, nil)
	MakeContextCurrent(w)
	for {
		PollEvents()
		SwapBuffers(w)
		if WindowShouldClose(w) != 0 {
			DestroyWindow(w)
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
