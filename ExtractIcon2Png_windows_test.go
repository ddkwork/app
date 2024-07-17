package app

import (
	"image/png"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/gorpher/gowin32"
)

func TestExtractIcon2Png(t *testing.T) {
	filename := "C:\\Program Files\\Tencent\\QQNT\\QQ.exe"
	img := mylog.Check2(gowin32.ExtractPrivateExtractIcons(filename, 128, 128))

	fp, _ := os.Create("output0.png")
	mylog.Check(png.Encode(fp, img))

	fp.Close()
}
