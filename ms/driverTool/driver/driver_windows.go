package driver

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ddkwork/golibrary/stream"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Load(deviceName, fileName string, Dependencies []string) error {
	log.Println("Loading Winpmem Driver...")

	content := stream.NewBuffer(fileName).Bytes()

	// write to file
	driverPath := filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", fileName)
	mylog.Check(os.WriteFile(driverPath, content, 0755))

	log.Println("Driver saved to", driverPath)

	// create service
	m := mylog.Check2(mgr.Connect())

	defer m.Disconnect()

	s, e := (m.OpenService(deviceName))
	if e == nil {
		s.Close()
		mylog.Check("serivce already exists")
	}
	config := mgr.Config{
		ServiceType:      windows.SERVICE_KERNEL_DRIVER,
		StartType:        mgr.StartManual,
		ErrorControl:     0,
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

	s = mylog.Check2(m.CreateService(deviceName, driverPath, config))

	defer s.Close()

	log.Println("Service created.")

	// start service
	mylog.Check(ControlService(deviceName, "start"))

	return nil
}

func Unload(deviceName string) error {
	log.Println("Unloading Winpmem Driver...")

	// Store driver to tempfile
	var driverName string
	if runtime.GOARCH == "386" {
		driverName = "winpmem_x86.sys"
	} else if runtime.GOARCH == "amd64" {
		driverName = "winpmem_x64.sys"
	} else {
		return errors.New("Architecture not supported: " + runtime.GOARCH)
	}

	driverPath := filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", driverName)

	// stop service
	mylog.Check(ControlService(deviceName, "stop"))

	// remove service
	mylog.Check(ControlService(deviceName, "delete"))
	// Delete driver file
	mylog.Check(os.Remove(driverPath))
	log.Printf("Drive file removed from: %v", driverPath)
	return nil
}

func ControlService(serviceName, action string) error {
	// open manager
	m := mylog.Check2(mgr.Connect())

	defer m.Disconnect()

	// open service
	s := mylog.Check2(m.OpenService(serviceName))

	defer s.Close()

	// stop service
	if action == "stop" {
		status := mylog.Check2(s.Control(svc.Stop))

		timeout := time.Now().Add(10 * time.Second)
		for status.State != svc.Stopped {
			if timeout.Before(time.Now()) {
				return errors.New("timed out waiting for service to stop")
			}
			time.Sleep(300 * time.Millisecond)
			status = mylog.Check2(s.Query())

		}
		log.Println("Service stopped.")
	}
	if action == "delete" {
		mylog.Check(s.Delete())
		log.Println("Service deleted.")
	}
	if action == "start" {
		mylog.Check(s.Start())
		log.Println("Service started")
	}

	return nil
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

func AcquireImage(deviceName, mode, filename string) error {
	Unload(deviceName)
	mylog.Check(Load(deviceName, filename, nil))
	defer Unload(deviceName)
	fd := mylog.Check2(syscall.CreateFile(
		syscall.StringToUTF16Ptr("\\\\.\\"+deviceName),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	))

	defer syscall.Close(fd)

	return nil
}
