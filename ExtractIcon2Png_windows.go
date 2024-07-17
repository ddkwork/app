package app

import (
	"bytes"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/gorpher/gowin32"
	"image/png"
)

func ExtractIcon2Png(filename string) []byte {
	img, e := gowin32.ExtractPrivateExtractIcons(filename, 128, 128)
	if e != nil {
		return nil
	}
	b := new(bytes.Buffer)
	mylog.Check(png.Encode(b, img))
	return b.Bytes()
}
