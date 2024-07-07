package glfw

import (
	"golang.org/x/sys/windows"
	"testing"
)

func StringToBytePointer(s string) *byte {
	bytes := []byte(s)
	ptr := &bytes[0]
	return ptr
}

func TestBind(t *testing.T) {
	windows.SetDllDirectory(".")
	Init()
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
