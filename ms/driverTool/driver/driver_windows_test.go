package driver

import (
	"testing"
)

func TestLoadSys(t *testing.T) {
	sysName := "sysDemo.sys"
	Load("", sysName, nil)
	Unload("", sysName)
}
