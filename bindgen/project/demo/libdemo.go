// Code generated by bindgen. DO NOT EDIT.
package libdemo

import (
	"unsafe"

	"github.com/ddkwork/app/bindgen/bindlib"
)

const GengoLibraryName = "libdemo"

var GengoLibrary = bindlib.NewLibrary(GengoLibraryName)

type Cr3Type struct {
	Anon274_5
}
type Anon274_5 struct {
	Raw [1]int64
}
type Anon278_9 struct {
	Pcid            Uint64
	PageFrameNumber Uint64
	Reserved1       Uint64
	Reserved_2      Uint64
	PcidInvalidate  Uint64
}
type (
	_Int128T           = any
	_Uint128T          = any
	__NSConstantString = any
	SizeT              = uint64
	_BuiltinMsVaList   = *byte
	_BuiltinVaList     = *byte
	Uint8T             = uint8
	Uint16T            = uint16
	Uint32T            = uint32
	Uint64T            = uint64
	Int8T              = int8
	Int16T             = int16
	Int32T             = int32
	Int64T             = int64
	Bool               = int32
	IntptrT            = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	// type Uint8T = uint8
	// type Uint16T = uint16
	// type Uint32T = uint32
	// type Uint64T = uint64
	// type Int8T = int8
	// type Int16T = int16
	// type Int32T = int32
	// type Int64T = int64
	// type Bool = int32
	// type IntptrT = *int32
	Uint64   = uint64
	Pcr3Type = *Cr3Type
)

var __imp_hello bindlib.PreloadProc

// Gengo init function.
func init() {
	__imp_hello = GengoLibrary.ImportNow("hello")
	bindlib.Validate((*Cr3Type)(nil), 8, 8)
	bindlib.Validate((*Anon274_5)(nil), 8, 8)
	bindlib.Validate((*Anon278_9)(nil), 8, 8, "Pcid", 0, "PageFrameNumber", 1, "Reserved1", 6, "Reserved_2", 7, "PcidInvalidate", 7)
}
func Hello() { bindlib.CCall0(__imp_hello.Addr()) }
func (s Anon274_5) Flags() Uint64 {
	return bindlib.ReadBitcast[Uint64](unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0))
}

func (s *Anon274_5) SetFlags(v Uint64) {
	bindlib.WriteBitcast(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0), v)
}

func (s Anon274_5) Fields() Anon278_9 {
	return bindlib.ReadBitcast[Anon278_9](unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0))
}

func (s *Anon274_5) SetFields(v Anon278_9) {
	bindlib.WriteBitcast(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0), v)
}
