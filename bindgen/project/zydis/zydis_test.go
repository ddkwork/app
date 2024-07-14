package zydis

import (
	"testing"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
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
	t.Skip()
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
	mylog.Check(pkg.WriteToDir("tmp"))
}
