package driver

import (
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestLoadSys(t *testing.T) {
	mylog.Call(func() {
		sysName := "sysDemo.sys"
		Load("", sysName)
		Unload("", sysName)
	})
}
