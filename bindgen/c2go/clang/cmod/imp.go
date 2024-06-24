package cmod

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ddkwork/app/bindgen/c2go/clang/pathutil"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/goplus/mod/gopmod"
	"github.com/qiniu/x/errors"
)

var ErrGoModNotFound = errors.New("go.mod not found")

func LoadDeps(dir string, deps []string) (pkgs []*Package, err error) {
	mod := mylog.Check2(gopmod.Load(dir))

	return Imports(mod, deps)
}

type Module = gopmod.Module

func Imports(mod *Module, pkgPaths []string) (pkgs []*Package, err error) {
	pkgs = make([]*Package, len(pkgPaths))
	for i, pkgPath := range pkgPaths {
		pkgs[i] = mylog.Check2(Import(mod, pkgPath))
	}
	return
}

type Package struct {
	*gopmod.Package
	Path    string   // package path
	Dir     string   // absolue local path of the package
	Include []string // absolute include paths
}

func Import(mod *Module, pkgPath string) (p *Package, err error) {
	if mod == nil {
		return nil, ErrGoModNotFound
	}
	pkg := mylog.Check2(mod.Lookup(pkgPath))

	pkgDir := mylog.Check2(filepath.Abs(pkg.Dir))

	pkgIncs := mylog.Check2(findIncludeDirs(pkgDir))

	for i, dir := range pkgIncs {
		pkgIncs[i] = pathutil.Canonical(pkgDir, dir)
	}
	return &Package{Package: pkg, Path: pkgPath, Dir: pkgDir, Include: pkgIncs}, nil
}

func findIncludeDirs(pkgDir string) (incs []string, err error) {
	var conf struct {
		Include []string `json:"include"`
	}
	file := filepath.Join(pkgDir, "c2go.cfg")
	b := mylog.Check2(os.ReadFile(file))
	mylog.Check(json.Unmarshal(b, &conf))

	return conf.Include, nil
}
