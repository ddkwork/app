package app

import (
	"os"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func TestUpdateAppModule(t *testing.T) {
	t.Skip()
	if !stream.IsDir("D:\\workspace\\workspace\\golibrary") {
		return
	}
	mylog.Check(os.Chdir("D:\\workspace\\workspace\\golibrary"))
	session := stream.RunCommand("git log -1 --format=\"%H\"")
	mylog.Check(os.Chdir("D:\\workspace\\workspace\\app"))
	id := mylog.Check2(strconv.Unquote(session.Output.String()))
	mylog.Info("id", id)
	stream.RunCommand("go get github.com/ddkwork/golibrary@" + id)
	//go get github.com/ddkwork/toolbox@e2865d103f34a0abdda09c5795343986799fb651
	//go get github.com/ddkwork/unison@b3a4eae98f926c8032f612c6e1a12622b6b72710
}
