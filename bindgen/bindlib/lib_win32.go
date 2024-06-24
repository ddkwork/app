//go:build windows

package bindlib

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
)

type windll struct {
	dll *windows.DLL
}

func (w windll) Lookup(name string) uintptr {
	proc := mylog.Check2(w.dll.FindProc(name))

	return proc.Addr()
}

func LoadLibrary(name string) (LoadedLibrary, error) {
	dll := mylog.Check2(windows.LoadDLL(name))

	return windll{dll: dll}, nil
}

func FindLibrary(name string) (LoadedLibrary, error) {
	return LoadLibrary(name)
}

func LoadLibraryEmbed(data []byte) (LoadedLibrary, error) {
	cache := getTmpDir()
	hash := sha1.Sum(data)
	name := "." + hex.EncodeToString(hash[:4]) + ".gengo.dll"
	path := cache + name
	if stat, e := os.Stat(path); e != nil || stat.Size() != int64(len(data)) {
		os.MkdirAll(cache, 0755)
		mylog.Check(os.WriteFile(path, data, 0755))

	}
	return LoadLibrary(path)
}
