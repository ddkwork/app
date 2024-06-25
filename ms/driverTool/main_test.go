package main

import (
	"testing"

	"github.com/ddkwork/app/ms/driverTool/driver"
)

func TestLoadSys(t *testing.T) {
	sysName := "sysDemo.sys"
	d := driver.New()
	d.Load(sysName)
	d.Unload()
}
