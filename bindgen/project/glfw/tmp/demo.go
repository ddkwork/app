package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/ddkwork/golibrary/stream"
	"golang.org/x/sys/windows"

	"github.com/ddkwork/golibrary/mylog"
)

//go:embed glfw3.dll
var dllData []byte

func init() {
	runtime.LockOSThread()
	dir := mylog.Check2(os.UserCacheDir())
	dir = filepath.Join(dir, "glfw3", "dll_cache")
	stream.CreatDirectory(dir)
	mylog.Check(windows.SetDllDirectory(dir))
	sha := sha256.Sum256(dllData)
	dllName := fmt.Sprintf("glfw3-%s.dll", base64.RawURLEncoding.EncodeToString(sha[:]))
	filePath := filepath.Join(dir, dllName)
	if !stream.IsFilePath(filePath) {
		stream.WriteTruncate(filePath, dllData)
	}
	mylog.Check2(GengoLibrary.LoadFrom(filePath))
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
