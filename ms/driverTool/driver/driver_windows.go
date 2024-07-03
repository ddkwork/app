package driver

import (
	"errors"
	"github.com/ddkwork/golibrary/stream"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Load(deviceName, fileName string, Dependencies []string) error {
	log.Println("Loading Winpmem Driver...")

	// Store driver to tempfile
	var driverName string
	if runtime.GOARCH == "386" {
		driverName = "winpmem_x86.sys"
	} else if runtime.GOARCH == "amd64" {
		driverName = "winpmem_x64.sys"
	} else {
		return errors.New("Architecture not supported: " + runtime.GOARCH)
	}

	content := stream.NewBuffer(fileName).Bytes()

	//write to file
	driverPath := filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", driverName)
	if err := os.WriteFile(driverPath, content, 0755); err != nil {
		return err
	}

	log.Println("Driver saved to", driverPath)

	//create service
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(deviceName)
	if err == nil {
		s.Close()
		return errors.New("serivce already exists")
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

	s, err = m.CreateService(deviceName, driverPath, config)
	if err != nil {
		return err
	}
	defer s.Close()

	log.Println("Service created.")

	//start service
	if err := ControlService(deviceName, "start"); err != nil {
		return err
	}

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

	//stop service
	if err := ControlService(deviceName, "stop"); err != nil {
		log.Printf("Unable to stop service: %v", err)
	}

	//remove service
	if err := ControlService(deviceName, "delete"); err != nil {
		log.Printf("Unable to delete service: %v", err)
	}

	//Delete driver file
	if err := os.Remove(driverPath); err != nil {
		log.Printf("Unable to remove driver file %v : %v", driverPath, err)
	} else {
		log.Printf("Drive file removed from: %v", driverPath)
	}
	return nil
}

func ControlService(serviceName, action string) error {
	//open manager
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	//open service
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	//stop service
	if action == "stop" {
		status, err := s.Control(svc.Stop)
		if err != nil {
			return err
		}
		timeout := time.Now().Add(10 * time.Second)
		for status.State != svc.Stopped {
			if timeout.Before(time.Now()) {
				return errors.New("timed out waiting for service to stop")
			}
			time.Sleep(300 * time.Millisecond)
			status, err = s.Query()
			if err != nil {
				return err
			}
		}
		log.Println("Service stopped.")
	}
	if action == "delete" {
		if err := s.Delete(); err != nil {
			log.Printf("Unable to delete service: %v", err)

		} else {
			log.Println("Service deleted.")
		}
	}
	if action == "start" {
		if err := s.Start(); err != nil {
			return err
		}
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
	if err := Load(deviceName, filename, nil); err != nil {
		return err
	}
	defer Unload(deviceName)
	fd, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr("\\\\.\\"+deviceName),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	return nil
}
