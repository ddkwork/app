package driver

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"runtime"
	"strings"
)

func CreateError(err error) error {
	if err == nil {
		return err
	}
	var callerName = "UnknownFunction"
	if info, _, _, ok := runtime.Caller(1); ok {
		details := runtime.FuncForPC(info)
		if details != nil {
			callerName = details.Name()
		}
	}
	callerNameSplit := strings.Split(callerName, ".")
	newErrorText := fmt.Sprintf("%s error: %s", callerNameSplit[len(callerNameSplit)-1], err.Error())
	return errors.New(newErrorText)
}

func GetProcessId(pid int, name string) int {
	if pid != 0 {
		return pid
	}
	processes, err := process.Processes()
	if err != nil {
		return 0
	}
	for _, each := range processes {
		if procName, err := each.Name(); err == nil && procName == name {
			return int(each.Pid)
		}
	}
	return 0
}
