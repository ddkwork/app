package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	logFile  = os.Stdout
	logLevel = LL_Info
)

type LogLevel uint

const (
	LL_Verbose LogLevel = iota
	LL_Debug
	LL_Info
	LL_Warn
	LL_Error
	LL_Critical
)

var logStrs = []string{"verbose", "debug", "info", "warn", "error", "critical"}

func SetLogFile(file string) {
	f := mylog.Check2(os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 644))

	logFile = f
}

func SetLogLevel(ll LogLevel) {
	if ll < LL_Verbose || ll > LL_Critical {
		return
	}
	logLevel = ll
}

func Log(ll LogLevel, format string, a ...any) {
	if ll < logLevel {
		return
	}
	if ll > LL_Critical {
		ll = LL_Critical
	}
	_, file, line, _ := runtime.Caller(1)
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	_, _ = fmt.Fprintf(logFile, "[%s] %s:%d %s: ", timeStr, filepath.Base(file), line, logStrs[ll])
	_, _ = fmt.Fprintf(logFile, format, a...)
	_, _ = fmt.Fprintln(logFile)
}

func Err(err error) {
	if err == nil {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	_, _ = fmt.Fprintf(logFile, "[%s] %s:%d %s: %s\n",
		timeStr, filepath.Base(file), line, logStrs[LL_Error], err.Error())
}
