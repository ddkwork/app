package main

import (
	"github.com/ddkwork/app"
	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/crypt/src/aes"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/orderedmap"
	"github.com/richardwilkes/unison"
)

func main() {
	app.Run("crypto tool", func(w *unison.Window) {
		content := w.Content()
		panel := widget.NewPanel()
		panel.AddChild(NewCryptUI().Layout())
		scrollPanelFill := widget.NewScrollPanelFill(panel)
		content.AddChild(scrollPanelFill)
	})
}

type CryptUI struct{}

func NewCryptUI() *CryptUI {
	return &CryptUI{}
}

type (
	CryptTable struct {
		Name CryptNameKind
	}
	SrcKeyDstdData struct {
		Src string
		Key string
		Dst string
	}
	HashData struct {
		Src    string
		Md2    string
		Md4    string
		Md5    string
		Sha1   string
		Sha224 string
		Sha256 string
		Sha384 string
		Sha512 string
	}
	RsaData struct {
		P string // 质数
		Q string // 质数
		E string // 指数
		N string // 模数
		D string // 私钥
		M string // 明文
		C string // 密文
	}
)

func (c *CryptUI) Layout() *unison.Panel {
	table, header := widget.NewTable(CryptTable{}, widget.TableContext[CryptTable]{
		ContextMenuItems: nil,
		MarshalRow: func(node *widget.Node[CryptTable]) (cells []widget.CellData) {
			name := node.Data.Name.String()
			if node.Container() {
				name = node.Sum(name)
			}
			return []widget.CellData{{Text: name}}
		},
		UnmarshalRow: func(node *widget.Node[CryptTable], values []string) {
			mylog.Todo("unmarshal row")
		},
		SelectionChangedCallback: func(root *widget.Node[CryptTable]) {
			mylog.Todo("selection changed callback")
		},
		SetRootRowsCallBack: func(root *widget.Node[CryptTable]) {
			for _, kind := range InvalidCryptKind.Kinds() {
				switch kind {
				case SymmetryKind:
					container := widget.NewContainerNode(SymmetryKind.String(), CryptTable{})
					root.AddChild(container)
					container.AddChildByData(CryptTable{Name: AesKind})
					container.AddChildByData(CryptTable{Name: DesKind})
					container.AddChildByData(CryptTable{Name: Des3Kind})
					container.AddChildByData(CryptTable{Name: TeaKind})
					container.AddChildByData(CryptTable{Name: BlowfishKind})
					container.AddChildByData(CryptTable{Name: TwoFishKind})
					container.AddChildByData(CryptTable{Name: Rc4Kind})
					container.AddChildByData(CryptTable{Name: Rc2Kind})
				case AsymmetricalKind:
					container := widget.NewContainerNode(AsymmetricalKind.String(), CryptTable{})
					root.AddChild(container)
					container.AddChildByData(CryptTable{Name: RsaKind})
					container.AddChildByData(CryptTable{Name: EccKind})
					container.AddChildByData(CryptTable{Name: DsaKind})
					container.AddChildByData(CryptTable{Name: PgpKind})
					container.AddChildByData(CryptTable{Name: Sm4Kind})
					container.AddChildByData(CryptTable{Name: Sm2Kind})
				case HashKind:
					container := widget.NewContainerNode(HashKind.String(), CryptTable{})
					root.AddChild(container)
					container.AddChildByData(CryptTable{Name: HmacKind})
					container.AddChildByData(CryptTable{Name: HashAllKind})
				case EncodingKind:
					container := widget.NewContainerNode(EncodingKind.String(), CryptTable{})
					root.AddChild(container)
					container.AddChildByData(CryptTable{Name: Base64Kind})
					container.AddChildByData(CryptTable{Name: Base32Kind})
					container.AddChildByData(CryptTable{Name: GzipKind})
				case ToolKind:
					container := widget.NewContainerNode(ToolKind.String(), CryptTable{})
					root.AddChild(container)
					container.AddChildByData(CryptTable{Name: TrimSpaceKind})
					container.AddChildByData(CryptTable{Name: SwapKind})
					container.AddChildByData(CryptTable{Name: RequestHeaderKind})
					container.AddChildByData(CryptTable{Name: TimeStampKind})
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
	layouts.Set(DesKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Des3Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(TeaKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(BlowfishKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(TwoFishKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Rc4Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Rc2Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(RsaKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(RsaData{}, func(data RsaData) (values []widget.CellData) {
			/*
				P string
				Q string
				E string
				N string
				D string
				M string
				C string

			*/
			return []widget.CellData{
				{Text: data.P},
				{Text: data.Q},
				{Text: data.E},
				{Text: data.N},
				{Text: data.D},
				{Text: data.M},
				{Text: data.C},
			}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
				// view.MetaData.Src = view.Editors[0].Label.String()
				// view.MetaData.Key = view.Editors[1].Label.String()
				// view.MetaData.Dst = string(aes.Encrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
				view.UpdateField(2, view.MetaData.C)
			},
			func() {
				// view.MetaData.Dst = view.Editors[2].Label.String()
				// view.MetaData.Key = view.Editors[1].Label.String()
				// view.MetaData.Src = string(aes.Decrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
				view.UpdateField(0, view.MetaData.M)
			},
		)
		RowPanel.AddChild(panel1)

		panel := widget.NewPanel()
		panel.AddChild(view)
		panel.AddChild(RowPanel)
		scrollPanelFill := widget.NewScrollPanelFill(panel)
		return scrollPanelFill
	})
	layouts.Set(EccKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(DsaKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(PgpKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Sm4Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Sm2Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(HmacKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(HashAllKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(HashData{}, func(data HashData) (values []widget.CellData) {
			return []widget.CellData{
				/*
							Src    string
					Md2    string
					Md4    string
					Md5    string
					Sha1   string
					Sha224 string
					Sha256 string
					Sha384 string
					Sha512 string

				*/
				{Text: data.Src},
				{Text: data.Md2},
				{Text: data.Md4},
				{Text: data.Md5},
				{Text: data.Sha1},
				{Text: data.Sha224},
				{Text: data.Sha256},
				{Text: data.Sha384},
				{Text: data.Sha512},
			}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
				view.MetaData.Src = view.Editors[0].Label.String()
				// view.MetaData.Key = view.Editors[1].Label.String()
				// view.MetaData.Dst = string(aes.Encrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
				// view.UpdateField(2, view.MetaData.Dst)
			},
			func() {
				// view.MetaData.Dst = view.Editors[2].Label.String()
				// view.MetaData.Key = view.Editors[1].Label.String()
				// view.MetaData.Src = string(aes.Decrypt(stream.HexString(view.MetaData.Src), stream.HexString(view.MetaData.Key)).HexString())
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
	layouts.Set(Base64Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Base32Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(GzipKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(TrimSpaceKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(SwapKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(RequestHeaderKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(TimeStampKind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
	layouts.Set(Base64Kind, func() unison.Paneler {
		view, RowPanel := widget.NewStructView(SrcKeyDstdData{}, func(data SrcKeyDstdData) (values []widget.CellData) {
			return []widget.CellData{{Text: data.Src}, {Text: data.Key}, {Text: data.Dst}}
		})
		panel1 := widget.NewButtonsPanel(
			[]string{"encode", "decode"},
			func() {
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
			case DesKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(DesKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Des3Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Des3Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case TeaKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(TeaKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case BlowfishKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(BlowfishKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case TwoFishKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(TwoFishKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Rc4Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Rc4Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Rc2Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Rc2Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case RsaKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(RsaKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case EccKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(EccKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case DsaKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(DsaKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case PgpKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(PgpKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Sm4Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Sm4Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Sm2Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Sm2Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case HmacKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(HmacKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case HashAllKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(HashAllKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case Base64Kind:
				right.RemoveAllChildren()
				panel := mylog.Check2Bool(layouts.Get(Base64Kind))()
				right.AddChild(panel)
				splitPanel.AddChild(right)
			case Base32Kind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(Base32Kind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case GzipKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(GzipKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case TrimSpaceKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(TrimSpaceKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case SwapKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(SwapKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case RequestHeaderKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(RequestHeaderKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			case TimeStampKind:
				right.RemoveAllChildren()
				paneler := mylog.Check2Bool(layouts.Get(TimeStampKind))()
				right.AddChild(paneler)
				splitPanel.AddChild(right)
			default:
			}
		}
	}
	return splitPanel.AsPanel()
}
