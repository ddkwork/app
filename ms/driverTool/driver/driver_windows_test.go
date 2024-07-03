package driver

import (
	"testing"
)

func TestLoadSys(t *testing.T) {
	path := "sysDemo.sys"
	Load("", path, nil)
	Unload("")
}
