package driver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ddkwork/golibrary/stream"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	ErrServiceStartPending = "SERVICE PENDING"
)

func Load(serviceName string, fileName string, Dependencies []string) {
	mylog.Call(func() {
		fileName = filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", filepath.Base(fileName))
		stream.WriteBinaryFile(fileName, stream.NewBuffer(fileName).Bytes())
		if serviceName == "" {
			serviceName = stream.BaseName(fileName)
		}
		mylog.Trace("deviceName", serviceName)
		mylog.Trace("driverPath", fileName)
		m := mylog.Check2(mgr.Connect())
		s := mylog.Check2Ignore(m.OpenService(serviceName))
		if s == nil {
			s = mylog.Check2(createService(m, serviceName, fileName, Dependencies))
		}
		if !verifyServiceConfig(s, fileName) {
			mylog.Check(m.Disconnect())
			mylog.Check(s.Close())
			Unload(serviceName, fileName)
			Load(serviceName, fileName, Dependencies)
		}
		mylog.Check(s.Start())
		verifyServiceRunning(serviceName)
		mylog.Success("driver load success", fileName)
	})
}

func Unload(serviceName string, fileName string) {
	mylog.Call(func() {
		m := mylog.Check2(mgr.Connect())
		service := mylog.Check2(m.OpenService(serviceName))
		if !verifyServiceConfig(service, fileName) {
			mylog.Check("invalid service")
		}
		mylog.Check2(service.Control(svc.Stop))
		mylog.Check(service.Delete())
		mylog.Success("driver unload success", fileName)
	})
}

func verifyServiceConfig(service *mgr.Service, driverPath string) bool {
	serviceConfig := mylog.Check2(service.Config())
	if serviceConfig.ServiceType != windows.SERVICE_KERNEL_DRIVER {
		return false
	}
	if serviceConfig.ErrorControl != windows.SERVICE_ERROR_IGNORE {
		return false
	}
	if serviceConfig.BinaryPathName != fmt.Sprintf("\\??\\%s", driverPath) {
		return false
	}
	return true
}

func verifyServiceRunning(serviceName string) {
	connSCM := mylog.Check2(mgr.Connect())
	service := mylog.Check2(connSCM.OpenService(serviceName))
	serviceStatus := mylog.Check2(service.Query())
	if serviceStatus.State == windows.SERVICE_START_PENDING {
		mylog.Check(ErrServiceStartPending)
	}
	if serviceStatus.State != windows.SERVICE_RUNNING {
		mylog.Check(errors.New("service was not started correctly"))
	}
}

func createService(m *mgr.Mgr, serviceName, driverPath string, Dependencies []string, args ...string) (*mgr.Service, error) {
	c := mgr.Config{
		ServiceType:      windows.SERVICE_KERNEL_DRIVER,
		StartType:        windows.SERVICE_DEMAND_START,
		ErrorControl:     windows.SERVICE_ERROR_IGNORE,
		BinaryPathName:   "",
		LoadOrderGroup:   "",
		TagId:            0,
		Dependencies:     Dependencies,
		ServiceStartName: "",
		DisplayName:      "",
		Password:         "",
		Description:      "",
		SidType:          0,
		DelayedAutoStart: false,
	}
	if c.StartType == 0 {
		c.StartType = mgr.StartManual
	}
	if c.ServiceType == 0 {
		c.ServiceType = windows.SERVICE_WIN32_OWN_PROCESS
	}
	h := mylog.Check2(windows.CreateService(m.Handle, toPtr(serviceName), toPtr(c.DisplayName),
		windows.SERVICE_ALL_ACCESS, c.ServiceType,
		c.StartType, c.ErrorControl, toPtr(driverPath), toPtr(c.LoadOrderGroup),
		nil, toStringBlock(c.Dependencies), toPtr(c.ServiceStartName), toPtr(c.Password)))

	if c.SidType != windows.SERVICE_SID_TYPE_NONE {
		updateSidType(h, c.SidType)
	}
	if c.Description != "" {
		updateDescription(h, c.Description)
	}
	if c.DelayedAutoStart {
		updateStartUp(h, c.DelayedAutoStart)
	}
	return &mgr.Service{Name: serviceName, Handle: h}, nil
}

func toPtr(s string) *uint16 {
	mylog.Check(len(s) == 0)
	return syscall.StringToUTF16Ptr(s)
}

func toStringBlock(ss []string) *uint16 {
	if len(ss) == 0 {
		return nil
	}
	t := ""
	for _, s := range ss {
		if s != "" {
			t += s + "\x00"
		}
	}
	if t == "" {
		return nil
	}
	t += "\x00"
	return &utf16.Encode([]rune(t))[0]
}

func updateSidType(handle windows.Handle, sidType uint32) {
	mylog.Check(windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_SERVICE_SID_INFO, (*byte)(unsafe.Pointer(&sidType))))
}

func updateDescription(handle windows.Handle, desc string) {
	d := windows.SERVICE_DESCRIPTION{Description: toPtr(desc)}
	mylog.Check(windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_DESCRIPTION, (*byte)(unsafe.Pointer(&d))))
}

func updateStartUp(handle windows.Handle, isDelayed bool) {
	var d windows.SERVICE_DELAYED_AUTO_START_INFO
	if isDelayed {
		d.IsDelayedAutoStartUp = 1
	}
	mylog.Check(windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_DELAYED_AUTO_START_INFO, (*byte)(unsafe.Pointer(&d))))
}

func GetProcessId(pid int, name string) int {
	if pid != 0 {
		return pid
	}
	processes := mylog.Check2(process.Processes())

	for _, each := range processes {
		if procName := mylog.Check2(each.Name()); procName == name {
			return int(each.Pid)
		}
	}
	return 0
}
