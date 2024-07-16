package zydis_test

import (
	"bytes"
	"fmt"
	"testing"
	"unsafe"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"

	"github.com/ddkwork/app/bindgen/project/zydis"
)

type zydisProvider struct {
	*gengo.BaseProvider
}

func (p *zydisProvider) NameField(name string, recordName string) string {
	if recordName == "ZydisDecodedInstructionRawEvex_" || recordName == "ZydisDecodedInstructionRawEvex" {
		if name == "b" {
			return "Br"
		}
	}
	return p.BaseProvider.NameField(name, recordName)
}

func TestZydis(t *testing.T) {
	prov := &zydisProvider{
		BaseProvider: gengo.NewBaseProvider(
			gengo.WithRemovePrefix(
				"Zydis_", "Zyan_", "Zycore_",
				"Zydis", "Zyan", "Zycore",
			),
			gengo.WithInferredMethods([]gengo.MethodInferenceRule{
				{Name: "ZydisDecoder", Receiver: "Decoder"},
				{Name: "ZydisEncoder", Receiver: "EncoderRequest"},
				{Name: "ZydisFormatterBuffer", Receiver: "FormatterBuffer"},
				{Name: "ZydisFormatter", Receiver: "ZydisFormatter *"},
				{Name: "ZyanVector", Receiver: "Vector"},
				{Name: "ZyanStringView", Receiver: "StringView"},
				{Name: "ZyanString", Receiver: "String"},
				{Name: "ZydisRegister", Receiver: "Register"},
				{Name: "ZydisMnemonic", Receiver: "Mnemonic"},
				{Name: "ZydisISASet", Receiver: "ISASet"},
				{Name: "ZydisISAExt", Receiver: "ISAExt"},
				{Name: "ZydisCategory", Receiver: "Category"},
			}),
			gengo.WithForcedSynthetic(
				"ZydisShortString_",
				"struct ZydisShortString_",
			),
		),
	}
	pkg := gengo.NewPackageWithProvider("zydis", prov)
	mylog.Check(pkg.Transform("zydis", &clang.Options{
		// Sources: []string{"codegen/Zydis.h"},
		// Sources: []string{"./Zydis.h"},
		Sources: []string{"Zydis.h"},
		AdditionalParams: []string{
			"-DZYAN_NO_LIBC",
			//"-DZYAN_STATIC_ASSERT",
		},
	}))
	mylog.Check(pkg.WriteToDir("."))
	Test_Disasm(t)
}

func Test_Disasm(t *testing.T) {
	fmt.Printf("Zydis Version: %x\n", zydis.GetVersion())

	data := []byte{
		0x51, 0x8D, 0x45, 0xFF, 0x50, 0xFF, 0x75, 0x0C, 0xFF, 0x75,
		0x08, 0xFF, 0x15, 0xA0, 0xA5, 0x48, 0x76, 0x85, 0xC0, 0x0F,
		0x88, 0xFC, 0xDA, 0x02, 0x00,
	}

	// The runtime address (instruction pointer) was chosen arbitrarily here in order to better
	// visualize relative addressing. In your actual program, set this to e.g. the memory address
	// that the code being disassembled was read from.
	runtimeAddress := uintptr(0x007FFFFFFF400000)

	// Loop over the instructions in our buffer.
	offset := 0
	insn := zydis.DisassembledInstruction{}
	for offset < len(data) {
		status := zydis.DisassembleIntel(
			zydis.MACHINE_MODE_LONG_64,
			uint64(runtimeAddress),
			unsafe.Pointer(&data[offset]),
			uint64(len(data)-offset),
			&insn,
		)
		if !zydis.Ok(status) {
			break
		}

		textEnd := bytes.IndexByte(insn.Text[:], 0)
		fmt.Printf("%016X  %s\n", runtimeAddress, string(insn.Text[:textEnd]))
		offset += int(insn.Info.Length)
		runtimeAddress += uintptr(insn.Info.Length)
	}
}
