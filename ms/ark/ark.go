package main

import (
	"github.com/ddkwork/app"
	"github.com/ddkwork/app/ms"
	"github.com/ddkwork/app/ms/driverTool/environment"
	"github.com/ddkwork/app/ms/hook/winver"
	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/orderedmap"
	"github.com/richardwilkes/unison"
)

func main() {
	app.Run("ark tool", func(w *unison.Window) {
		content := w.Content()
		panel := widget.NewPanel()
		panel.AddChild(Layout())
		scrollPanelFill := widget.NewScrollPanelFill(panel)
		content.AddChild(scrollPanelFill)
	})
}

func arkTodo() {
	ms.DecodeTableByDll()
	println(winver.WindowVersion())
	ms.MiGetPteAddress()
	ms.DecodeTableByDll()
	ms.DecodeTableByDisassembly()
	ms.NtDeviceIoControlFile()
	// IopXxxControlFile()
	widget.NewExplorer(".")
	environment.New()
}

func Layout() *unison.Panel {
	type ark struct{ Name ArksKind }
	table, header := widget.NewTable(ark{}, widget.TableContext[ark]{
		ContextMenuItems: nil,
		MarshalRow: func(node *widget.Node[ark]) (cells []widget.CellData) {
			name := node.Data.Name.String()
			if node.Container() {
				name = node.Sum(name)
			}
			return []widget.CellData{{Text: name}}
		},
		UnmarshalRow: func(node *widget.Node[ark], values []string) {
			mylog.Todo("unmarshal row")
		},
		SelectionChangedCallback: func(root *widget.Node[ark]) {
			mylog.Todo("selection changed callback")
		},
		SetRootRowsCallBack: func(root *widget.Node[ark]) {
			for _, kind := range InvalidArksKind.Kinds() {
				root.AddChildByData(ark{kind})
			}
		},
		JsonName:   "ark",
		IsDocument: false,
	})

	splitPanel := widget.NewPanel()
	widget.SetScrollLayout(splitPanel, 2)

	left := widget.NewTableScrollPanel(table, header)
	layouts := orderedmap.New(InvalidArksKind, func() unison.Paneler { return widget.NewPanel() })
	layouts.Set(KernelTablesKind, func() unison.Paneler {
		return widget.NewButton("111", nil)
	})

	right := widget.NewPanel()
	right.AddChild(mylog.Check2Bool(layouts.Get(KernelTablesKind))()) // todo make a welcoming page
	splitPanel.AddChild(left)
	splitPanel.AddChild(right)

	// todo get and set inputted ctx,not clean it every time
	table.SelectionChangedCallback = func() {
		for i, n := range table.SelectedRows(false) {
			if i > 1 {
				break
			}
			switch n.Data.Name {
			case KernelTablesKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(KernelTablesKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)

			default:
			}
		}
	}
	return splitPanel.AsPanel()
}
