package translate

import (
	"clang/ast"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ddkwork/golibrary/mylog"
)

type transErr struct {
	*ast.Loc
	msg  string
	file string
	line int
}

func newTransErr(loc *ast.Loc, v ...any) error {
	err := &transErr{Loc: loc}
	err.msg = fmt.Sprint(v...)
	_, err.file, err.line, _ = runtime.Caller(2)
	err.file = filepath.Base(err.file)
	return err
}

func (s *transErr) Error() string {
	mylog.Check(fmt.Sprintf("%s:%d %s:%d:%d %s", s.file, s.line, s.File, s.Line, s.Col, s.msg))
	from := s.IncludedFrom
	for from != nil {
		err += "\n\tfrom: " + from.File
		from = from.IncludedFrom
	}
	return err
}
