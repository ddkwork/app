package main

import (
	"testing"

	"github.com/goradd/maps"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	g := stream.NewGeneratedFile()
	m := new(maps.SafeSliceMap[string, string])
	m.Set("smoke", "烟")
	m.Set("wine", "酒")
	m.Set("meat", "肉")
	m.Set("disk", "菜")
	m.Set("sheep", "羊")
	m.Set("gift", "礼金")
	m.Set("cash", "取钱")
	m.Set("firecrackers", "鞭炮")
	m.Set("other", "其它")
	g.EnumTypes("trade", m)
}
