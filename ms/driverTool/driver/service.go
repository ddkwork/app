package driver

import (
	"errors"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	ErrServiceStartPending = "SERVICE PENDING"
)

func SetUpService(serviceName string, driverFullPath string) {
	m := mylog.Check2(mgr.Connect())
	s := mylog.Check2Ignore(m.OpenService(serviceName))
	if s == nil {
		s = mylog.Check2(CreateService(m, serviceName, driverFullPath))
	}
	if !VerifyServiceConfig(s, driverFullPath) {
		mylog.Check(m.Disconnect())
		mylog.Check(s.Close())
		RemoveService(serviceName, driverFullPath)
		SetUpService(serviceName, driverFullPath)
	}
	mylog.Check(s.Start())
}

func VerifyServiceConfig(service *mgr.Service, driverPath string) bool {
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

func VerifyServiceRunning(serviceName string) error {
	connSCM := mylog.Check2(mgr.Connect())
	service := mylog.Check2(connSCM.OpenService(serviceName))
	serviceStatus := mylog.Check2(service.Query())
	if serviceStatus.State == windows.SERVICE_START_PENDING {
		return errors.New(ErrServiceStartPending)
	} else if serviceStatus.State != windows.SERVICE_RUNNING {
		mylog.Check(errors.New("service was not started correctly"))
	}
	return nil
}

func RemoveService(serviceName string, driverFullPath string) {
	m := mylog.Check2(mgr.Connect())
	service := mylog.Check2(m.OpenService(serviceName))
	if !VerifyServiceConfig(service, driverFullPath) {
		mylog.Check(errors.New("invalid service"))
	}
	mylog.Check2(service.Control(svc.Stop))
	mylog.Check(service.Delete())
}

func CreateService(m *mgr.Mgr, serviceName, driverPath string, args ...string) (*mgr.Service, error) {
	c := mgr.Config{
		ServiceType:  windows.SERVICE_KERNEL_DRIVER,
		StartType:    windows.SERVICE_DEMAND_START,
		ErrorControl: windows.SERVICE_ERROR_IGNORE,
	}
	if c.StartType == 0 {
		c.StartType = mgr.StartManual
	}
	if c.ServiceType == 0 {
		c.ServiceType = windows.SERVICE_WIN32_OWN_PROCESS
	}
	h := mylog.Check2(windows.CreateService(m.Handle, toUnicode(serviceName), toUnicode(c.DisplayName),
		windows.SERVICE_ALL_ACCESS, c.ServiceType,
		c.StartType, c.ErrorControl, toUnicode(driverPath), toUnicode(c.LoadOrderGroup),
		nil, toStringBlock(c.Dependencies), toUnicode(c.ServiceStartName), toUnicode(c.Password)))

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

func toUnicode(s string) *uint16 {
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
	d := windows.SERVICE_DESCRIPTION{Description: toUnicode(desc)}
	mylog.Check(windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_DESCRIPTION, (*byte)(unsafe.Pointer(&d))))
}

func updateStartUp(handle windows.Handle, isDelayed bool) {
	var d windows.SERVICE_DELAYED_AUTO_START_INFO
	if isDelayed {
		d.IsDelayedAutoStartUp = 1
	}
	mylog.Check(windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_DELAYED_AUTO_START_INFO, (*byte)(unsafe.Pointer(&d))))
}
