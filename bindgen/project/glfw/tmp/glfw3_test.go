package glfw3_test

import (
	glfw3 "github.com/ddkwork/app/bindgen/project/glfw/tmp"
	"golang.org/x/sys/windows"
	"testing"
)

func TestBind(t *testing.T) {
	windows.SetDllDirectory(".")
	glfw3.Init()
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
