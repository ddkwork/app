package packer

import (
	"strings"

	"github.com/ddkwork/app/ms/xed"

	"github.com/ddkwork/golibrary/mylog"
)

func DumpPe() {
	checkPacker("D:\\clone\\demo\\unlicense\\unlicense-py3.11-x86\\SuperRecovery.exe")
}

func checkPacker(filename string) {
	file := xed.ParserPe(filename)
	for _, section := range file.Sections {
		s := section.String()
		after, found := strings.CutPrefix(s, ".")
		if !found {
			continue
		}
		s = after
		for _, v := range SigMap {
			if strings.Contains(v, s) {
				mylog.Warning("packet", s)
			}
		}
	}
}
