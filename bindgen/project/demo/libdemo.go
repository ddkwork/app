// Code generated by bindgen. DO NOT EDIT.
package libdemo

import (
	"github.com/ddkwork/app/bindgen/bindlib"
	"unsafe"
)

const GengoLibraryName = "libdemo"

var GengoLibrary = bindlib.NewLibrary(GengoLibraryName)

type Cr3Type struct {
	Anon286_5
}
type Anon286_5 struct {
	Raw [1]int64
}
type Anon290_9 struct {
	Pcid            Uint64
	PageFrameNumber Uint64
	Reserved1       Uint64
	Reserved_2      Uint64
	PcidInvalidate  Uint64
}
type _Int128T = any
type _Uint128T = any
type __NSConstantString = any
type SizeT = uint64
type _BuiltinMsVaList = *byte
type _BuiltinVaList = *byte
type Uint8T = uint8
type Uint16T = uint16
type Uint32T = uint32
type Uint64T = uint64
type Int8T = int8
type Int16T = int16
type Int32T = int32
type Int64T = int64
type Bool = int32
type IntptrT = *int32

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
type Uint64 = uint64
type Pcr3Type = *Cr3Type

var __imp_hello bindlib.PreloadProc

// Gengo init function.
func init() {
	__imp_hello = GengoLibrary.ImportNow("hello")
	bindlib.Validate((*Cr3Type)(nil), 8, 8)
	bindlib.Validate((*Anon286_5)(nil), 8, 8)
	bindlib.Validate((*Anon290_9)(nil), 8, 8, "Pcid", 0, "PageFrameNumber", 1, "Reserved1", 6, "Reserved_2", 7, "PcidInvalidate", 7)
}
func Hello() { bindlib.CCall0(__imp_hello.Addr()) }
func (s Anon286_5) Flags() Uint64 {
	return bindlib.ReadBitcast[Uint64](unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0))
}
func (s *Anon286_5) SetFlags(v Uint64) {
	bindlib.WriteBitcast(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0), v)
}
func (s Anon286_5) Fields() Anon290_9 {
	return bindlib.ReadBitcast[Anon290_9](unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0))
}
func (s *Anon286_5) SetFields(v Anon290_9) {
	bindlib.WriteBitcast(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s.Raw[:])), 0), v)
}
