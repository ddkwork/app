package main

import (
	"github.com/ddkwork/app/ms/hardwareIndo"
)

func test() {
	h := hardwareIndo.New()
	if !h.SsdInfo.Get() { // todo bug cpu pkg init
		return
	}
	if !h.CpuInfo.Get() {
		return
	}
	if !h.MacInfo.Get() {
		return
	}
}
