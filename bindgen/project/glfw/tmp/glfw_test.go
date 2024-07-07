package glfw

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestInit(t *testing.T) {
	windows.SetDllDirectory(".")
	Init()
	//CreateWindow(200,200)
	Terminate()
	// GetError()
	// SwapBuffers(nil)
}
