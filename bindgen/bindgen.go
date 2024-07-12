package main

type (
	Binder interface {
		BindCgo()
		BindSharedLibrary()
		C2go()   // TODO qiniu cc
		CPP2go() // TODO gitee cc
		// skia ndk
	}
)

// clang -E -P temp.h > output_file.txt
// https://git.homegu.com/KhronosGroup/OpenGL-Registry/tree/main/api/GLES3
