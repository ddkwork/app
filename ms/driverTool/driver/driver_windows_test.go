package driver

import (
	"testing"
)

func TestLoadSys(t *testing.T) {
	path := "sysDemo.sys"
	d := New("", path, nil)
	d.Load()
	d.Unload()
}
