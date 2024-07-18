package main

import "github.com/ddkwork/app/ms"

//go:generate go build -x .
//go:generate go install .
func main() {
	ms.DecodeTableByDll()
}
