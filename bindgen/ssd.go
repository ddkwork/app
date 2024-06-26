package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var HardDriveSerialNumber [1024]byte

const (
	DFP_RECEIVE_DRIVE_DATA = 0x0007c088

	FILE_DEVICE_DISK                = 0x00000007
	FILE_DEVICE_SCSI                = 0x0000001b
	IOCTL_SCSI_MINIPORT_IDENTIFY    = ((FILE_DEVICE_SCSI << 16) + 0x0501)
	IOCTL_SCSI_MINIPORT             = 0x0004D008
	SMART_GET_VERSION               = 0x00074080
	SMART_RCV_DRIVE_DATA            = 0x0007c088
	IOCTL_STORAGE_QUERY_PROPERTY    = 0x002d1400
	IOCTL_STORAGE_GET_DEVICE_NUMBER = 0x002d1080
	IDENTIFY_BUFFER_SIZE            = 512
)

type IDSECTOR struct {
	WGenConfig                 uint16
	WNumCyls                   uint16
	WReserved                  uint16
	WNumHeads                  uint16
	WBytesPerTrack             uint16
	WBytesPerSector            uint16
	WSectorsPerTrack           uint16
	WVendorUnique              [3]uint16
	SSerialNumber              [20]byte
	WBufferType                uint16
	WBufferSize                uint16
	WECCSize                   uint16
	SFirmwareRev               [8]byte
	SModelNumber               [40]byte
	WMoreVendorUnique          uint16
	WDoubleWordIO              uint16
	WCapabilities              uint16
	WReserved1                 uint16
	WPIOTiming                 uint16
	WDMATiming                 uint16
	WBS                        uint16
	WNumCurrentCyls            uint16
	WNumCurrentHeads           uint16
	WNumCurrentSectorsPerTrack uint16
	UlCurrentSectorCapacity    uint32
	WMultSectorStuff           uint16
	UlTotalAddressableSectors  uint32
	WSingleWordDMA             uint16
	WMultiWordDMA              uint16
	BReserved                  [128]byte
}

type SRB_IO_CONTROL struct {
	HeaderLength uint32
	Signature    [8]byte
	Timeout      uint32
	ControlCode  uint32
	ReturnCode   uint32
	Length       uint32
}

type StoragePropertyQuery struct {
	PropertyId int32
	QueryType  int32
}

type StorageDeviceDescriptor struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        byte
	CommandQueueing       byte
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               byte
	RawPropertiesLength   uint32
	RawDeviceProperties   [1024]byte
}

func ReadPhysicalDriveInNTUsingSmart(drive int) bool {
	done := false
	var driveName = fmt.Sprintf("\\\\.\\PhysicalDrive%d", drive)
	driveNameUTF16, _ := syscall.UTF16PtrFromString(driveName)

	hPhysicalDriveIOCTL, err := syscall.CreateFile(driveNameUTF16, syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		fmt.Printf(" ReadPhysicalDriveInNTUsingSmart ERROR: CreateFile(%s) returned %d \n\n", driveName, err)
	} else {
		defer syscall.CloseHandle(hPhysicalDriveIOCTL)
		var GetVersionParams [IDENTIFY_BUFFER_SIZE]byte
		var cbBytesReturned uint32

		if err := syscall.DeviceIoControl(hPhysicalDriveIOCTL, SMART_GET_VERSION,
			nil, 0,
			(*byte)(unsafe.Pointer(&GetVersionParams[0])), uint32(len(GetVersionParams)),
			&cbBytesReturned, nil); err != nil {
			panic(err)
		} else {
			const ID_CMD = 0xEC
			CommandSize := int(unsafe.Sizeof(SRB_IO_CONTROL{})) + IDENTIFY_BUFFER_SIZE
			Command := make([]byte, CommandSize)
			Command[0] = ID_CMD

			if err := syscall.DeviceIoControl(hPhysicalDriveIOCTL,
				SMART_RCV_DRIVE_DATA, (*byte)(unsafe.Pointer(&Command[0])), uint32(CommandSize), (*byte)(unsafe.Pointer(&Command[0])), uint32(CommandSize),
				&cbBytesReturned, nil); err != nil {
				fmt.Printf("SMART_RCV_DRIVE_DATA IOCTL error: %d\n", err)
			} else {
				var diskdata [256]uint16
				pIdSector := (*IDSECTOR)(unsafe.Pointer(&Command[0]))

				for i := 0; i < 256; i++ {
					diskdata[i] = *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(pIdSector)) + uintptr(i*2)))
				}

				PrintIdeInfo(drive, diskdata[:])
				done = true
			}
		}
	}
	return done
}

func ConvertToString(diskdata []uint16, firstIndex, lastIndex int, buf *[]byte) {
	for i := firstIndex; i <= lastIndex; i++ {
		*buf = append(*buf, byte(diskdata[i]>>8), byte(diskdata[i]&0xFF))
	}
	index := 0
	for index < len(*buf) && (*buf)[index] != 0 {
		if (*buf)[index] == ' ' {
			copy((*buf)[index:], (*buf)[index+1:])
			*buf = (*buf)[:len(*buf)-1]
		} else {
			index++
		}
	}
}

func PrintIdeInfo(drive int, diskdata []uint16) {
	ConvertToString(diskdata, 10, 19, (*[]byte)(unsafe.Pointer(&HardDriveSerialNumber)))

	fmt.Printf("\nDrive %d - ", drive)

	switch drive / 2 {
	case 0:
		fmt.Println("Primary Controller - ")
	case 1:
		fmt.Println("Secondary Controller - ")
	case 2:
		fmt.Println("Tertiary Controller - ")
	case 3:
		fmt.Println("Quaternary Controller - ")
	}

	switch drive % 2 {
	case 0:
		fmt.Println("Master drive\n\n")
	case 1:
		fmt.Println("Slave drive\n\n")
	}

	fmt.Printf("Drive Serial Number_______________: [%s]\n", strings.TrimSpace(string(HardDriveSerialNumber[:])))
	fmt.Printf("Drive Type________________________: ")

	if diskdata[0]&0x0080 != 0 {
		fmt.Println("Removable")
	} else if diskdata[0]&0x0040 != 0 {
		fmt.Println("Fixed")
	} else {
		fmt.Println("Unknown")
	}
}

func SystemDrive() int {
	drive := 0
	hPart, err := os.Open("\\\\.\\C:")
	if err != nil {
		fmt.Printf("Open ERROR %d\n\n", err)
	} else {
		defer hPart.Close()

		var Info [512]byte
		var tmp uint32

		if err := syscall.DeviceIoControl(syscall.Handle(hPart.Fd()), IOCTL_STORAGE_GET_DEVICE_NUMBER, nil, 0,
			(*byte)(unsafe.Pointer(&Info[0])), uint32(len(Info)), &tmp, nil); err != nil {
			fmt.Printf("DeviceIoControl ERROR %d\n\n", err)
		} else {
			ptr := unsafe.Pointer(&Info[0])
			deviceType := *(*uint32)(ptr)
			deviceNumber := *(*int)(unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(uint32(0))))

			if deviceType == FILE_DEVICE_DISK {
				drive = deviceNumber
				fmt.Printf("System drive: %d\n\n", drive)
			} else {
				fmt.Println("Info missized or bad DeviceType\n\n")
			}
		}
	}

	return drive
}

func main() {
	//drive := SystemDrive()
	fmt.Printf("\nTrying to read the drive IDs using physical access with zero rights\n\n")
	ReadPhysicalDriveInNTUsingSmart(0)

	fmt.Printf("\nHard Drive Serial Number__________: %s\n\n", HardDriveSerialNumber[:])
}
