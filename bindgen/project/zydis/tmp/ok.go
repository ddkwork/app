package zydis

import "golang.org/x/sys/windows"

func Ok(status uint32) bool {
	return int32(status) >= 0
}

//func Failed(status Status) bool {
//	return int32(status) < 0
//}

func init() {
	windows.SetDllDirectory(".")
}
