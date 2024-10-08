/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package c2go

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ddkwork/app/bindgen/c2go/cl"
	"github.com/ddkwork/golibrary/mylog"

	"github.com/ddkwork/app/bindgen/c2go/clang/parser"
	"github.com/ddkwork/app/bindgen/c2go/clang/preprocessor"
)

const (
	FlagRunApp = 1 << iota
	FlagRunTest
	FlagFailFast
	FlagDepsAutoGen
	FlagForcePreprocess
	FlagDumpJson
	FlagTestMain

	flagChdir
)

func isDir(name string) bool {
	if fi, e := (os.Lstat(name)); e == nil {
		return fi.IsDir()
	}
	return false
}

func isFile(name string) bool {
	if fi, e := (os.Lstat(name)); e == nil {
		return !fi.IsDir()
	}
	return false
}

type Config struct {
	SelectFile string
	SelectCmd  string
}

func Run(pkgname, infile string, flags int, conf *Config) {
	outfile := infile
	switch filepath.Ext(infile) {
	case ".i":
	case ".c":
		outfile = infile + ".i"
		mylog.Check(preprocessor.Do(infile, outfile, nil))
	default:
		if strings.HasSuffix(infile, "/...") {
			infile = strings.TrimSuffix(infile, "/...")
			mylog.Check(execDirRecursively(infile, flags, conf))
		} else if isDir(infile) {
			projfile := filepath.Join(infile, "c2go.cfg")
			if isFile(projfile) {
				execProj(projfile, flags, conf)
				return
			}
			n := mylog.Check2(execDir(pkgname, infile, flags))
			switch n {
			case 1:
			case 0:
				fatalf("no *.c files in this directory.\n")
			default:
				fatalf("multiple .c files found (currently only support one .c file).\n")
			}
		} else {
			fatalf("%s is not a .c file.\n", infile)
		}
		return
	}
	execFile(pkgname, outfile, flags|flagChdir)
	return
}

func execDirRecursively(dir string, flags int, conf *Config) (last error) {
	if strings.HasPrefix(dir, "_") {
		return
	}

	projfile := filepath.Join(dir, "c2go.cfg")
	if isFile(projfile) {
		fmt.Printf("==> Compiling %s ...\n", dir)
		execProj(projfile, flags, conf)
		return
	}

	fis := mylog.Check2(os.ReadDir(dir))
	var cfiles int
	for _, fi := range fis {
		if fi.IsDir() {
			pkgDir := filepath.Join(dir, fi.Name())
			if e := execDirRecursively(pkgDir, flags, conf); e != nil {
				last = e
			}
			continue
		}
		if strings.HasSuffix(fi.Name(), ".c") {
			cfiles++
		}
	}
	if cfiles == 1 {
		var action string
		switch {
		case (flags & FlagRunTest) != 0:
			action = "Testing"
		case (flags & FlagRunApp) != 0:
			action = "Running"
		default:
			action = "Compiling"
		}
		fmt.Printf("==> %s %s ...\n", action, dir)
		if _, e := execDir("main", dir, flags); e != nil {
			last = e
		}
	}
	return
}

func execDir(pkgname string, dir string, flags int) (n int, err error) {
	if (flags & FlagFailFast) == 0 {
		defer func() {
			if e := recover(); e != nil {
				mylog.Check(newError(e))
			}
		}()
	}
	n = -1

	cwd := chdir(dir)
	defer os.Chdir(cwd)

	var infile, outfile string
	files := mylog.Check2(filepath.Glob("*.c"))
	switch n = len(files); n {
	case 1:
		infile = files[0]
		outfile = infile + ".i"
		mylog.Check(preprocessor.Do(infile, outfile, nil))
		execFile(pkgname, outfile, flags)
	}
	return
}

func execFile(pkgname string, outfile string, flags int) {
	var json []byte
	doc, _ := mylog.Check3(parser.ParseFileEx(outfile, 0, &parser.Config{
		Json:   &json,
		Stderr: true,
	}))

	if (flags & FlagDumpJson) != 0 {
		os.WriteFile(strings.TrimSuffix(outfile, ".i")+".json", json, 0666)
	}

	needPkgInfo := (flags & FlagDepsAutoGen) != 0
	pkg := mylog.Check2(cl.NewPackage("", pkgname, doc, &cl.Config{
		SrcFile: outfile, NeedPkgInfo: needPkgInfo,
	}))

	gofile := outfile + ".go"
	mylog.Check(pkg.WriteFile(gofile))

	dir, _ := filepath.Split(gofile)

	if needPkgInfo {
		mylog.Check(pkg.WriteDepFile(filepath.Join(dir, "c2go_autogen.go")))
	}

	if (flags & flagChdir) != 0 {
		if dir != "" {
			cwd := chdir(dir)
			defer os.Chdir(cwd)
		}
	}

	if (flags & FlagRunTest) != 0 {
		runTest("")
	} else if (flags & FlagRunApp) != 0 {
		runGoApp("", os.Stdout, os.Stderr, false)
	}
}

func checkEqual(prompt string, a, expected []byte) {
	if bytes.Equal(a, expected) {
		return
	}

	fmt.Fprintln(os.Stderr, "=> Result of", prompt)
	os.Stderr.Write(a)

	fmt.Fprintln(os.Stderr, "\n=> Expected", prompt)
	os.Stderr.Write(expected)

	fatal(errors.New("checkEqual: unexpected " + prompt))
}

func runTest(dir string) {
	var goOut, goErr bytes.Buffer
	var cOut, cErr bytes.Buffer
	dontRunTest := runGoApp(dir, &goOut, &goErr, true)
	if dontRunTest {
		return
	}
	runCApp(dir, &cOut, &cErr)
	checkEqual("output", goOut.Bytes(), cOut.Bytes())
	checkEqual("stderr", goErr.Bytes(), cErr.Bytes())
}

func runGoApp(dir string, stdout, stderr io.Writer, doRunTest bool) (dontRunTest bool) {
	files := mylog.Check2(filepath.Glob("*.go"))

	for i, n := 0, len(files); i < n; i++ {
		fname := filepath.Base(files[i])
		if pos := strings.LastIndex(fname, "_"); pos >= 0 {
			switch os := fname[pos+1 : len(fname)-3]; os {
			case "darwin", "linux", "windows":
				if os != runtime.GOOS { // skip
					n--
					files[i], files[n] = files[n], files[i]
					files = files[:n]
					i--
				}
			}
		}
	}

	if doRunTest {
		for _, file := range files {
			if filepath.Base(file) == "main.go" {
				stdout, stderr = os.Stdout, os.Stderr
				dontRunTest = true
				break
			}
		}
	}
	cmd := exec.Command("go", append([]string{"run"}, files...)...)
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	checkWith(cmd.Run(), stdout, stderr)
	return
}

func runCApp(dir string, stdout, stderr io.Writer) {
	files := mylog.Check2(filepath.Glob("*.c"))

	cmd := exec.Command("clang", files...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	mylog.Check(cmd.Run())

	cmd2 := exec.Command(clangOut)
	cmd.Dir = dir
	cmd2.Stdout = stdout
	cmd2.Stderr = stderr
	checkWith(cmd2.Run(), stdout, stderr)

	os.Remove(clangOut)
}

var clangOut = "./a.out"

func init() {
	if runtime.GOOS == "windows" {
		clangOut = "./a.exe"
	}
}

func chdir(dir string) string {
	cwd := mylog.Check2(os.Getwd())
	mylog.Check(os.Chdir(dir))
	return cwd
}

func checkWith(err error, stdout, stderr io.Writer) {
}

func fatalf(format string, args ...interface{}) {
	fatal(fmt.Errorf(format, args...))
}

func fatal(err error) {
	log.Panicln(err)
}

func fatalWith(err error, stdout, stderr io.Writer) {
	if o, ok := getBytes(stdout, stderr); ok {
		os.Stderr.Write(o.Bytes())
	}
	log.Panicln(err)
}

func newError(v interface{}) error {
	switch e := v.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	}
	fatalf("newError failed: %v", v)
	return nil
}

type iBytes interface {
	Bytes() []byte
}

func getBytes(stdout, stderr io.Writer) (o iBytes, ok bool) {
	if o, ok = stderr.(iBytes); ok {
		return
	}
	o, ok = stdout.(iBytes)
	return
}
