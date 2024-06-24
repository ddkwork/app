package cmd

import (
	_ "embed"
	"testing"
)

//go:embed tmp/1.c
var cFile []byte

//go:embed tmp/merged_headers.h
var merged_headers []byte

//go:embed vmm/ept/Ept.c
var ceptFile []byte

func Test_translate(t *testing.T) {
	dumpAST = true
	translate(cFile)
}
