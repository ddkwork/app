package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ddkwork/app/bindgen/c2go/clang/parser"
)

var dump = flag.Bool("dump", false, "dump AST")

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: clangast [-dump] source.i\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		return
	}
	file := flag.Arg(0)

	if *dump {
		doc, _, e := parser.DumpAST(file, nil)
		if e == nil {
			os.Stdout.Write(doc)
			return
		}
		err = e
	} else {
		doc, _, e := parser.ParseFile(file, 0)
		if e == nil {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(doc)
			return
		}
		err = e
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
