package main

import (
	"github.com/ddkwork/golibrary/stream"
	"testing"
)

func TestName(t *testing.T) {
	g := stream.NewGeneratedFile()
	g.Types(
		"trade", //交易类型
		[]string{
			"smoke",
			"wine",
			"meat",
			"disk",
			"sheep",
			"gift",
			"cash",
			"firecrackers",
			"other",
		},
		[]string{
			"烟",
			"酒",
			"肉",
			"菜",
			"羊",
			"礼金",
			"取钱",
			"鞭炮",
			"其它",
		},
	)
}
