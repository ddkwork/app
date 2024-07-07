package glfw

import (
	"fmt"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestInit(t *testing.T) {
	windows.SetDllDirectory(".")
	Init()
	//CreateWindow(200,200)
	DestroyWindow(nil)
	Terminate()
	// GetError()
	// SwapBuffers(nil)
}

func main() {
	// Load the GLFW DLL
	glfw := syscall.MustLoadDLL("glfw3.dll")

	// Get the function pointers
	glfwInit := glfw.MustFindProc("glfwInit")
	glfwCreateWindow := glfw.MustFindProc("glfwCreateWindow")
	glfwMakeContextCurrent := glfw.MustFindProc("glfwMakeContextCurrent")
	glfwPollEvents := glfw.MustFindProc("glfwPollEvents")
	glfwTerminate := glfw.MustFindProc("glfwTerminate")

	// Initialize GLFW
	ret, _, err := glfwInit.Call()
	if ret == 0 {
		fmt.Println("Failed to initialize GLFW:", err)
		return
	}
	defer glfwTerminate.Call()

	// Create a windowed mode window and its OpenGL context
	width, height := 640, 480
	title := "Hello World"
	window, _, err := glfwCreateWindow.Call(
		uintptr(width), uintptr(height), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	if window == 0 {
		fmt.Println("Failed to create GLFW window:", err)
		return
	}

	// Make the window's context current
	glfwMakeContextCurrent.Call(window)

	// Main loop
	for {
		// Poll for and process events
		glfwPollEvents.Call()

		// Swap buffers
		// (You would typically call your rendering code here)

		// Check if the window should close
		shouldClose, _, _ := syscall.Syscall(glfw.MustFindProc("glfwWindowShouldClose").Addr(), 1, window, 0, 0)
		if shouldClose != 0 {
			break
		}
	}
}
