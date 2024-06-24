package project

import (
	"go/scanner"
	"go/token"

	"github.com/ddkwork/golibrary/stream"
)

func ScanComments(path string) (ok bool) {
	b := stream.NewBuffer(path)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), b.Len())
	s.Init(file, b.Bytes(), nil /* no error handler */, scanner.ScanComments)

	for {
		_, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		switch {
		case tok == token.COMMENT:
		default:
		}
	}
	return true
}
