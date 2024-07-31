package skia

import (
	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestMergeHeader(t *testing.T) {
	b := stream.NewBuffer("c/sk_types.h")
	filepath.Walk("c", func(path string, info fs.FileInfo, err error) error {
		if filepath.Base(path) != "sk_types.h" {
			b.WriteStringLn(stream.NewBuffer(path).String())
		}
		return nil
	})
	stream.WriteTruncate("skia.h", b.Bytes())
}

func TestBindSkia(t *testing.T) {
	//TestMergeHeader(t)
	pkg := gengo.NewPackage("skia")
	path := "skia.h"
	mylog.Check(pkg.Transform("skia", &clang.Options{
		Sources: []string{path},
		//AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("tmp"))
}
