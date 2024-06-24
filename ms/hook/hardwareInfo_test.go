package hook

import (
	"testing"

	"github.com/ddkwork/app/ms/hardwareIndo"

	"github.com/ddkwork/golibrary/mylog"
)

func Test_hardware(t *testing.T) {
	t.Skip()
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
	mylog.Struct(h)
}
