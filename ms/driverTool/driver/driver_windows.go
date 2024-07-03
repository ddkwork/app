package driver

import (
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr" // todo if build on linux,it need change to cmd

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type (
	Object struct {
		Status       uint32
		service      *mgr.Service
		manager      *mgr.Mgr
		path         string
		DeviceName   string
		Dependencies []string
	}
)

func New(deviceName, driverPath string, Dependencies []string) (d *Object) {
	return &Object{
		Status:       0,
		service:      nil,
		manager:      nil,
		path:         driverPath,
		DeviceName:   deviceName,
		Dependencies: Dependencies,
	}
}

func (o *Object) Load() {
	if o.path == "" {
		o.path = filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", filepath.Base(o.path))
		stream.WriteBinaryFile(o.path, stream.NewBuffer(o.path).Bytes())
	}
	if o.DeviceName == "" {
		o.DeviceName = stream.BaseName(o.path)
	}
	mylog.Trace("deviceName", o.DeviceName)
	mylog.Trace("path", o.path)
	o.manager = mylog.Check2(mgr.Connect())
	o.SetService()
	o.StartService()
	mylog.Success("driver load success", o.path)
	o.QueryService()
}

func (o *Object) Unload() {
	o.StopService()
	o.DeleteService()
	mylog.Check(o.manager.Disconnect())
	mylog.Check(o.service.Close())
	mylog.Success("driver unload success", o.path)
	mylog.Check(os.Remove(o.path))
}

func (o *Object) SetService() {
	var e error
	o.service, e = o.manager.OpenService(o.DeviceName)
	if e != nil {
		config := mgr.Config{
			ServiceType:      windows.SERVICE_KERNEL_DRIVER,
			StartType:        mgr.StartManual,
			ErrorControl:     0,
			BinaryPathName:   "",
			LoadOrderGroup:   "",
			TagId:            0,
			Dependencies:     o.Dependencies,
			ServiceStartName: "",
			DisplayName:      "",
			Password:         "",
			Description:      "",
			SidType:          0,
			DelayedAutoStart: false,
		}
		o.service = mylog.Check2(o.manager.CreateService(o.DeviceName, o.path, config))
	}
}

func (o *Object) QueryService() {
	o.Status = mylog.Check2(o.service.Query()).ServiceSpecificExitCode
}

func (o *Object) StopService() {
	status := mylog.Check2(o.service.Control(svc.Stop))
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			mylog.Check("Timed out waiting for service to stop")
		}
		time.Sleep(300 * time.Millisecond)
		o.QueryService()
		mylog.Trace("Service stopped")
	}
}

func (o *Object) DeleteService() {
	mylog.Check(o.service.Delete())
	mylog.Trace("Service deleted")
	o.QueryService()
}
func (o *Object) StartService() { mylog.Check(o.service.Start()) }
