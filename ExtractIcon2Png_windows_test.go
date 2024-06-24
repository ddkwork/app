package app

import (
	"image/png"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestExtractIcon2Png(t *testing.T) {
	return
	path := "D:\\app\\Internet Download Manager 6.42 Build 1 多语+Retail 坡姐\\Internet Download Manager 6.42 Build 1 多语+Retail 坡姐\\Patch-Ali.Dbg_v18.2\\IDM v.6.4x crack v.18.2.exe"
	path = "C:\\Windows\\notepad.exe"
	mylog.Call(func() {
		image, ok := ExtractIcon2Image(path)
		mylog.Check(ok)
		f := mylog.Check2(os.Create("1.png"))
		mylog.Check(png.Encode(f, image))
		mylog.Check(f.Close())
		mylog.Check(os.Remove("1.png"))
	})
}
