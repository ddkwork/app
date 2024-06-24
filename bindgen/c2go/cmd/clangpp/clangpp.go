package main

import (
	"fmt"
	"os"

	"github.com/ddkwork/app/bindgen/c2go/clang/preprocessor"
	"github.com/ddkwork/golibrary/mylog"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: clangpp source.c\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	infile := os.Args[1]
	outfile := infile + ".i"
	if mylog.Check(preprocessor.Do(infile, outfile, nil)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
