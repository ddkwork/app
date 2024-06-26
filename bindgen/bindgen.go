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
