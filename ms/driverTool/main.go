package main

import (
	_ "embed"
	"io/fs"
	"path/filepath"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/app/ms/driverTool/driver"

	"github.com/ddkwork/app"
	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison"
)

//go:embed icon.png
var icon []byte

func main() {
	app.Run("driver tool", func(w *unison.Window) {
		img := mylog.Check2(unison.NewImageFromBytes(icon, 0.5))
		w.SetTitleIcons([]*unison.Image{img})
		w.Content().AddChild(New().Layout())
	})
}

type DriverLoad struct {
	ReloadPath string
	Link       string
	IoCode     string
}

type StructView struct{}

func New() widget.API {
	return &StructView{}
}

func (s *StructView) Layout() unison.Paneler {
	view := DriverLoad{
		ReloadPath: "",
		Link:       "",
		IoCode:     "",
	}
	structView, rowPanel := widget.NewStructView(view, func(data DriverLoad) (values []widget.CellData) {
		return []widget.CellData{
			{ImageBuffer: nil, Text: data.ReloadPath, Tooltip: "", FgColor: 0},
			{ImageBuffer: nil, Text: data.Link, Tooltip: "", FgColor: 0},
			{ImageBuffer: nil, Text: data.IoCode, Tooltip: "", FgColor: 0},
		}
	})

	p := unison.NewPopupMenu[string]()

	p.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
		if title, ok := popup.Selected(); ok {
			structView.MetaData.ReloadPath = title
			structView.MetaData.Link = stream.BaseName(title)
			structView.UpdateField(0, title)
			structView.UpdateField(1, stream.BaseName(title))
		}
	}

	root := "../"
	abs := mylog.Check2(filepath.Abs(root))
	names := make([]string, 0)
	if abs == "C:\\Users\\Admin" {
		names = WalkAllDriverPath(".")
	} else {
		names = WalkAllDriverPath(root)
	}

	popupMenu := widget.CreatePopupMenu(rowPanel, p, 0, "choose a driver", names...)

	kv := widget.NewKeyValuePanel()
	key := widget.NewLabelRightAlign(widget.KeyValueToolTip{
		Key:     "sys path",
		Value:   "",
		Tooltip: "",
	})
	kv.AddChild(key)
	kv.AddChild(popupMenu)
	structView.AddChildAtIndex(kv, 0) // todo bug need rowPanel AddChildAtIndex

	d := driver.NewObject()
	log := unison.NewMultiLineField() // todo log out format is not good
	log.MinimumTextWidth = 800
	log.SetText(`log view






`)
	panel := widget.NewButtonsPanel(
		[]string{"load", "unload"},
		func() {
			d.Load(structView.MetaData.ReloadPath)
			log.SetText(mylog.Body())
		},
		func() {
			d.Unload()
			log.SetText(mylog.Body())
		},
	)
	structView.AddChild(widget.NewVSpacer())
	structView.AddChild(panel)
	structView.AddChild(log)
	return structView
}

func WalkAllDriverPath(root string) (drivers []string) {
	drivers = make([]string, 0)
	mylog.Check(filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".sys" {
			drivers = append(drivers, path)
		}
		return nil
	}))
	return
}
