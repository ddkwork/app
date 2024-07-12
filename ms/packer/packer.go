package packer

import (
	"strings"

	"github.com/ddkwork/app/ms/xed"

	"github.com/ddkwork/golibrary/mylog"
)

func CheckPacker(filename string) {
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
