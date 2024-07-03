package driver

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/shirou/gopsutil/process"
)

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
