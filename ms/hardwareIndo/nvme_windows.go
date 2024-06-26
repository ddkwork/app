package hardwareIndo

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
)

// 定义系统常量
const (
	IOCTL_STORAGE_QUERY_PROPERTY  = 0x2D1400
	IOCTL_DISK_GET_DRIVE_GEOMETRY = 0x70000
)

// 定义结构
type HDiskInfo struct {
	module   []byte // 40
	firmware [8]byte
	serialno [20]byte
	capacity uint32
}

type STORAGE_PROPERTY_QUERY struct {
	PropertyId           uint32
	QueryType            uint32
	AdditionalParameters [1]byte
}

type STORAGE_DEVICE_DESCRIPTOR struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        bool
	CommandQueueing       bool
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               byte
	RawPropertiesLength   uint32
}

type DISK_GEOMETRY struct {
	Cylinders         int64
	MediaType         uint32
	TracksPerCylinder uint32
	SectorsPerTrack   uint32
	BytesPerSector    uint32
}

func readHarddiskInfo(pinfo *HDiskInfo) int {
	if pinfo == nil {
		return -1
	}

	// 打开物理驱动器
	hDevice := mylog.Check2(syscall.CreateFile(
		syscall.StringToUTF16Ptr(`\\.\PhysicalDrive0`),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	))

	defer syscall.CloseHandle(hDevice)

	query := STORAGE_PROPERTY_QUERY{
		PropertyId: 0,
		QueryType:  0,
	}

	buffer := make([]byte, 1024)

	var bytesReturned uint32
	mylog.
		// 查询硬盘属性
		Check(syscall.DeviceIoControl(
			hDevice,
			IOCTL_STORAGE_QUERY_PROPERTY,
			(*byte)(unsafe.Pointer(&query)),
			uint32(unsafe.Sizeof(query)),
			&buffer[0],
			uint32(len(buffer)),
			&bytesReturned,
			nil,
		))

	deviceDescriptor := (*STORAGE_DEVICE_DESCRIPTOR)(unsafe.Pointer(&buffer[0]))

	if deviceDescriptor.SerialNumberOffset != 0 {
		copy(pinfo.serialno[:], buffer[deviceDescriptor.SerialNumberOffset:])
	}
	if deviceDescriptor.ProductIdOffset != 0 {
		all := buffer[deviceDescriptor.ProductIdOffset:]
		end := 0
		size := 0
		for i, b := range all {
			if b == 0 {
				end = i + int(deviceDescriptor.ProductIdOffset)
				all = buffer[deviceDescriptor.ProductIdOffset:end]
				size = i
				break
			}
		}
		pinfo.module = make([]byte, size)
		copy(pinfo.module[:size], all)
	}
	if deviceDescriptor.ProductRevisionOffset != 0 {
		copy(pinfo.firmware[:], buffer[deviceDescriptor.ProductRevisionOffset:])
	}

	// 查询硬盘容量
	geom := DISK_GEOMETRY{}
	mylog.Check(syscall.DeviceIoControl(
		hDevice,
		IOCTL_DISK_GET_DRIVE_GEOMETRY,
		nil,
		0,
		(*byte)(unsafe.Pointer(&geom)),
		uint32(unsafe.Sizeof(geom)),
		&bytesReturned,
		nil,
	))
	pinfo.capacity = uint32(geom.Cylinders) * geom.TracksPerCylinder * geom.SectorsPerTrack * geom.BytesPerSector / (1024 * 1024)
	return 0
}

func nvme() {
	var hddInfo HDiskInfo
	status := readHarddiskInfo(&hddInfo)

	if status == 0 {
		fmt.Printf("硬盘信息获取成功:\n")
		fmt.Printf("型号: %s\n", hddInfo.module)
		fmt.Printf("固件版本: %s\n", hddInfo.firmware)
		fmt.Printf("序列号: %s\n", hddInfo.serialno)
		fmt.Printf("容量: %d MB\n", hddInfo.capacity)
	} else {
		fmt.Printf("硬盘信息获取失败，错误码: %d\n", status)
	}
}
