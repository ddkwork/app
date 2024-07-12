package main

import (
	"fmt"

	"github.com/ddkwork/app"
	"github.com/ddkwork/app/widget"
	"github.com/richardwilkes/unison"
)

func main() {
	app.Run("cpuInfo", func(w *unison.Window) {
		w.Content().AddChild(LayoutCpuInfo())
	})
}

func LayoutCpuInfo() unison.Paneler {
	type (
		Data0 struct {
			Arg int
			EAX int
			EBX int
			ECX int
			EDX int
		}
		Data1 struct {
			Arg int
			EAX int
			EBX int
			ECX int
			EDX int
		}
	)

	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{Columns: 2})

	view0, kvPanel0 := widget.NewStructView(Data0{}, func(data Data0) (values []widget.CellData) {
		return []widget.CellData{
			{Text: fmt.Sprintf("%016X", data.Arg)},
			{Text: fmt.Sprintf("%016X", data.EAX)},
			{Text: fmt.Sprintf("%016X", data.EBX)},
			{Text: fmt.Sprintf("%016X", data.ECX)},
			{Text: fmt.Sprintf("%016X", data.EDX)},
		}
	})
	view1, kvPanel1 := widget.NewStructView(Data0{}, func(data Data0) (values []widget.CellData) {
		return []widget.CellData{
			{Text: fmt.Sprintf("%016X", data.Arg)},
			{Text: fmt.Sprintf("%016X", data.EAX)},
			{Text: fmt.Sprintf("%016X", data.EBX)},
			{Text: fmt.Sprintf("%016X", data.ECX)},
			{Text: fmt.Sprintf("%016X", data.EDX)},
		}
	})
	panel.AddChild(view0)
	panel.AddChild(kvPanel0)
	panel.AddChild(view1)
	panel.AddChild(kvPanel1)
	return panel
}
