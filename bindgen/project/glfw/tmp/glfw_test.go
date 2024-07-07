package glfw

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestInit(t *testing.T) {
	windows.SetDllDirectory(".")
	Init()
	// GetError()
	// SwapBuffers(nil)
}
