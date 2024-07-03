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
	connSCM := mylog.Check2(mgr.Connect())
	service := CheckService(*connSCM, serviceName)
	if service == nil {
		service = mylog.Check2(CreateService(connSCM, serviceName, driverFullPath))
	}
	if !VerifyServiceConfig(service, driverFullPath) {
		mylog.Check(connSCM.Disconnect())
		mylog.Check(service.Close())
		RemoveService(serviceName, driverFullPath)
		SetUpService(serviceName, driverFullPath)
	}
	mylog.Check(service.Start())
}

func CheckService(connSCM mgr.Mgr, serviceName string) *mgr.Service {
	return mylog.Check2(connSCM.OpenService(serviceName))
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
	connSCM := mylog.Check2(mgr.Connect())
	service := mylog.Check2(connSCM.OpenService(serviceName))
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
		mylog.Check(updateSidType(h, c.SidType))
	}
	if c.Description != "" {
		mylog.Check(updateDescription(h, c.Description))
	}
	if c.DelayedAutoStart {
		mylog.Check(updateStartUp(h, c.DelayedAutoStart))
	}
	return &mgr.Service{Name: serviceName, Handle: h}, nil
}

func toUnicode(s string) *uint16 {
	if len(s) == 0 {
		return nil
	}
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

func updateSidType(handle windows.Handle, sidType uint32) error {
	return windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_SERVICE_SID_INFO, (*byte)(unsafe.Pointer(&sidType)))
}

func updateDescription(handle windows.Handle, desc string) error {
	d := windows.SERVICE_DESCRIPTION{Description: toUnicode(desc)}
	return windows.ChangeServiceConfig2(handle,
		windows.SERVICE_CONFIG_DESCRIPTION, (*byte)(unsafe.Pointer(&d)))
}

func updateStartUp(handle windows.Handle, isDelayed bool) error {
	var d windows.SERVICE_DELAYED_AUTO_START_INFO
	if isDelayed {
		d.IsDelayedAutoStartUp = 1
	}
	return windows.ChangeServiceConfig2(handle, windows.SERVICE_CONFIG_DELAYED_AUTO_START_INFO, (*byte)(unsafe.Pointer(&d)))
}
