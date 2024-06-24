package xed

import (
	"github.com/saferwall/pe"

	"github.com/ddkwork/golibrary/mylog"
)

func ParserPe(filename string) (file *pe.File) {
	file = mylog.Check2(pe.New(filename, &pe.Options{}))
	mylog.Check(file.Parse())
	return file
}
