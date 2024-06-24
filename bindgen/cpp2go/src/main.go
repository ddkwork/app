package main

import (
	"clang/parser"
	"io/fs"
	"path/filepath"
	"strings"
	"translate"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

func transFile(cSrcFile, goSrcFile string) {
	conf := &parser.Config{
		// Json:       "test/ast.json",
		Flags:  []string{"-fno-builtin", "-nostdinc"},
		Stderr: true,
		// SkipOthers: true,
	}
	file := mylog.Check2(parser.ParseFile(cSrcFile, conf))
	// file, err := parser.LoadJsonFile("test/ast.Json")

	utils.Log(utils.LL_Info, "==> start to translate")
	cfg := translate.Config{Package: "main", FileName: cSrcFile, SkipOthers: true}
	trans := translate.NewTranslate(file, &cfg)
	utils.Err(trans.Run())
	utils.Err(trans.Save(goSrcFile))
	utils.Log(utils.LL_Info, "==> end to translate")
}

func transProject(cSrcPath, goSrcPath, pkg string) {
	conf := &parser.Config{
		Flags:      []string{"-fno-builtin", "-nostdinc"},
		Stderr:     true,
		SkipOthers: true,
	}
	mylog.Check(filepath.WalkDir(cSrcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		name := d.Name()
		if !isCSrcFile(name) {
			return nil
		}
		file := mylog.Check2(parser.ParseFile(path, conf))

		if index := strings.LastIndex(name, "."); index > 0 {
			name = name[:index]
		}
		goFile := filepath.Join(goSrcPath, name+".go")
		utils.Log(utils.LL_Info, "==> start to translate %s -> %s", path, goFile)
		cfg := translate.Config{Package: pkg, FileName: path, SkipOthers: true}
		trans := translate.NewTranslate(file, &cfg)
		utils.Err(trans.Run())
		utils.Err(trans.Save(goFile))
		utils.Log(utils.LL_Info, "==> end to translate %s", path)
		return nil
	}))
	utils.Err(err)
}

func main() {
	utils.SetLogLevel(utils.LL_Info)
	// transFile("test/hello.c", "test/out.go")
	transProject("test", "test", "main")
}

func isCSrcFile(fileName string) bool {
	cExtName := []string{".c", ".cc", ".cpp"}
	for _, ext := range cExtName {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
