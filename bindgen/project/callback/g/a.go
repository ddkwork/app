package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// 定义回调函数类型
type FileDropCallbackFunc func(files []string)

// 全局变量存储回调函数
var dragHandler FileDropCallbackFunc

// 设置回调函数
func FileDropCallback(fn FileDropCallbackFunc) {
	dragHandler = fn
}

func main() {
	// 加载 DLL
	dllPath := "callback.dll"
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		fmt.Println("加载 DLL 失败:", err)
		return
	}
	defer dll.Release()

	// 获取 SetFileDropCallback 函数
	setCallbackProc, err := dll.FindProc("SetFileDropCallback")
	if err != nil {
		fmt.Println("查找 SetFileDropCallback 函数失败:", err)
		return
	}

	// 设置回调函数
	callback := syscall.NewCallback(func(files uintptr, fileCount int) uintptr {
		fileList := make([]string, fileCount)
		for i := 0; i < fileCount; i++ {
			filePtr := *(*uintptr)(unsafe.Pointer(files + uintptr(i)*unsafe.Sizeof(uintptr(0))))
			fileList[i] = syscall.UTF16ToString((*[1 << 30]uint16)(unsafe.Pointer(filePtr))[:])
		}
		dragHandler(fileList)
		return 0
	})

	// 调用 SetFileDropCallback 函数并传递回调函数
	_, _, err = setCallbackProc.Call(callback)
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("调用 SetFileDropCallback 函数失败:", err)
		return
	}

	// 设置回调函数
	FileDropCallback(func(files []string) {
		fmt.Println("拖放的文件:", files)
	})

	// 获取 TriggerCallback 函数
	triggerCallbackProc, err := dll.FindProc("TriggerCallback")
	if err != nil {
		fmt.Println("查找 TriggerCallback 函数失败:", err)
		return
	}

	// 准备文件列表
	fileList := []string{"file1.txt", "file2.txt"}
	filePointers := make([]uintptr, len(fileList))
	for i, file := range fileList {
		filePointers[i] = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(file)))
	}

	// 调用 TriggerCallback 函数
	_, _, err = triggerCallbackProc.Call(uintptr(unsafe.Pointer(&filePointers[0])), uintptr(len(fileList)))
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("调用 TriggerCallback 函数失败:", err)
		return
	}
}
