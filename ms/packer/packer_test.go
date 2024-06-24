package packer

import (
	"strings"
	"testing"
)

// D:\clone\VMProtect-devirtualization-main
// D:\workspace\hv\unlicense\unlicense-py3.11-x86
func TestDumpPe(t *testing.T) {
	DumpPe()
}

func TestDecodeSysCall(t *testing.T) {
	// println(strings.Contains(".themida\n", "themida-winlicense"))
	after, found := strings.CutPrefix(".themida", ".")
	if !found {
		return
	}
	println(strings.Contains("themida-winlicense", after))
}
