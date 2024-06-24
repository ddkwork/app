package main

import "github.com/ddkwork/app/ms"

//go:generate go build .
//go:generate go install .
func main() {
	ms.DecodeTableByDll()
}
