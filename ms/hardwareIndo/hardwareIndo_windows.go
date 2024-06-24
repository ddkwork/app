package hardwareIndo

type (
	Interface interface{}
	object    struct {
		MacInfo *systemScanner
		CpuInfo *cpuInfo
		SsdInfo *ssdInfo
	}
)

func New() *object {
	return &object{
		MacInfo: new(systemScanner),
		CpuInfo: new(cpuInfo),
		SsdInfo: new(ssdInfo),
	}
}
