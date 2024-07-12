package main

import (
	"github.com/ddkwork/app"
	"github.com/ddkwork/app/ms"
	"github.com/ddkwork/app/ms/driverTool/environment"
	"github.com/ddkwork/app/ms/hook/winver"
	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/crypt/src/aes"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
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
	type ark struct {
		KernelTables      string
		Explorer          string
		TaskManager       string
		DriverTool        string
		RegistryEditor    string
		HardwareMonitor   string
		HardwareHook      string
		RandomHook        string
		EnvironmentEditor string
		Vstart            string
		Crypt             string
		InvalidArks       string
	}

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
			for _, kind := range InvalidCryptKind.Kinds() {
				switch kind {
				case SymmetryKind:
					container := widget.NewContainerNode(SymmetryKind.String(), ark{})
					root.AddChild(container)
					container.AddChildByData(ark{Name: AesKind})
					container.AddChildByData(ark{Name: DesKind})
					container.AddChildByData(ark{Name: Des3Kind})
					container.AddChildByData(ark{Name: TeaKind})
					container.AddChildByData(ark{Name: BlowfishKind})
					container.AddChildByData(ark{Name: TwoFishKind})
					container.AddChildByData(ark{Name: Rc4Kind})
					container.AddChildByData(ark{Name: Rc2Kind})

				default:
				}
			}
		},
		JsonName:   "Crypt",
		IsDocument: false,
	})

	splitPanel := widget.NewPanel()
	widget.SetScrollLayout(splitPanel, 2)

	left := widget.NewTableScrollPanel(table, header)
	layouts := orderedmap.New(InvalidCryptNameKind, func() unison.Paneler { return widget.NewPanel() })
	layouts.Set(AesKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
				if view.Editors[0].Label.String() == "" { // todo
					view.Editors[0].Label.SetTitle("1122334455667788")
				}
				if view.Editors[1].Label.String() == "" {
					view.Editors[1].Label.SetTitle("1122334455667788")
				}

				view.MetaData.Src = view.Editors[0].Label.String()
				view.MetaData.Key = view.Editors[1].Label.String()
				view.MetaData.Dst = string(aes.Encrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
				view.UpdateField(2, view.MetaData.Dst)
			},
			func() {
				view.MetaData.Dst = view.Editors[2].Label.String()
				view.MetaData.Key = view.Editors[1].Label.String()
				view.MetaData.Src = string(aes.Decrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
				view.UpdateField(0, view.MetaData.Src)
			},
		)
		RowPanel.AddChild(panel1)

		panel := widget.NewPanel()
		panel.AddChild(view)
		panel.AddChild(RowPanel)
		scrollPanelFill := widget.NewScrollPanelFill(panel)
		return scrollPanelFill
	})

	right := widget.NewPanel()
	right.AddChild(mylog.Check2Bool(layouts.Get(AesKind))()) // todo make a welcoming page
	splitPanel.AddChild(left)
	splitPanel.AddChild(right)

	// todo get and set inputted ctx,not clean it every time
	table.SelectionChangedCallback = func() {
		for i, n := range table.SelectedRows(false) {
			if i > 1 {
				break
			}
			switch n.Data.Name {
			case AesKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(AesKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)

			default:
			}
		}
	}
	return splitPanel.AsPanel()
}
