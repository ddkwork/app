//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc
// +build 386 arm armbe mips mipsle ppc s390 sparc

package clang

type (
	Long  = int32
	Ulong = uint32
)
