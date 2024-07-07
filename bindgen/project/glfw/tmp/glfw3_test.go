package glfw

import (
	"syscall"
	"testing"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
)

func StringToBytePointer(s string) *byte {
	bytes := []byte(s)
	ptr := &bytes[0]
	return ptr
}

func TestInit(t *testing.T) {
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
	ret, _ := mylog.Check3(glfwInit.Call())
	if ret == 0 {
		// fmt.Println("Failed to initialize GLFW:", err)
		return
	}
	defer glfwTerminate.Call()

	// Create a windowed mode window and its OpenGL context
	width, height := 640, 480
	title := "Hello World"
	window, _ := mylog.Check3(glfwCreateWindow.Call(
		uintptr(width), uintptr(height), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0))
	if window == 0 {
		// fmt.Println("Failed to create GLFW window:", err)
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
