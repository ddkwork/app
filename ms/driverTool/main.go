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
		New(w).Layout(w.Content())
	})
}

type DriverLoad struct {
	ReloadPath string
	Link       string
	IoCode     string
}

type StructView struct {
	w *unison.Window
}

func New(w *unison.Window) widget.API {
	return &StructView{w: w}
}

func (s *StructView) Layout(parent unison.Paneler) unison.Paneler {
	view := DriverLoad{
		ReloadPath: "",
		Link:       "",
		IoCode:     "",
	}
	structView, rowPanel := widget.NewStructView(s.w, view, func(data DriverLoad) (values []widget.CellData) {
		return []widget.CellData{
			{ImageBuffer: nil, Text: data.ReloadPath, Tooltip: "", FgColor: 0},
			{ImageBuffer: nil, Text: data.Link, Tooltip: "", FgColor: 0},
			{ImageBuffer: nil, Text: data.IoCode, Tooltip: "", FgColor: 0},
		}
	})
	widget.NewLabelKey(widget.KeyValueToolTip{
		Parent:  rowPanel,
		Key:     "driver name",
		Value:   "",
		Tooltip: "",
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
	rowPanel.AddChild(popupMenu)

	d := driver.NewObject()
	log := widget.NewMultiLineField()
	log.MinimumTextWidth = 520
	log.SetText("loading...")
	rowPanel.AddChild(log)
	widget.NewButtonsPanel(structView,
		[]string{"load", "unload"},
		func() {
			d.Load(structView.MetaData.ReloadPath)
			// log.SetBuf(texteditor.NewBuf().SetText([]byte(mylog.Reason())))// todo
			log.SetText(structView.MetaData.ReloadPath + " load successfully.")
		},
		func() {
			d.Unload()
			// log.SetBuf(texteditor.NewBuf().SetText([]byte(mylog.Reason())))// todo
			log.SetText(structView.MetaData.ReloadPath + " unload successfully.")
		},
	)
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
