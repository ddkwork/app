package widget

import (
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
)

type CodeEditor struct {
	unison.Panel
	Editor   *Field
	undoMgr  *unison.UndoManager
	scroller *unison.ScrollPanel
}

func NewCodeEditor(parent unison.Paneler, filePath string) *CodeEditor {
	c := &CodeEditor{
		Panel:    unison.Panel{},
		Editor:   NewMultiLineField(),
		undoMgr:  unison.NewUndoManager(200, func(err error) {}),
		scroller: unison.NewScrollPanel(),
	}
	c.Self = c
	c.SetLayout(&unison.FlexLayout{Columns: 1})
	c.scroller.SetLayoutData(&unison.FlexLayoutData{
		SizeHint: unison.Size{},
		MinSize:  unison.Size{},
		HSpan:    0,
		VSpan:    0,
		HAlign:   align.Fill,
		VAlign:   align.Fill,
		HGrab:    true,
		VGrab:    true,
	})
	c.Editor.InitialClickSelectsAll = func(field *Field) bool {
		return false
	}
	c.scroller.SetContent(c.Editor, behavior.Fill, behavior.Fill)
	c.AddChild(c.scroller)

	c.Editor.SetText(filePath)
	//c.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
	//	c.scroller.SetPosition(where.X, where.Y)
	//	c.scroller.ScrollIntoView()
	//	return true
	//}
	c.Editor.AsPanel().FileDropCallback = func(files []string) {
		switch filepath.Ext(files[0]) {
		case ".go", ".scala":
			c.Editor.SetText(files[0])
		default:
			mylog.Check("file not go or scala")
		}
	}

	NewContextMenuItems(c.Editor, c.Editor.DefaultMouseDown,
		ContextMenuItem{
			Title: "Copy",
			Can:   func(any) bool { return c.Editor.CanCopy() },
			Do:    func(a any) { c.Editor.Copy() },
		},
		ContextMenuItem{
			Title: "Paste",
			Can:   func(any) bool { return c.Editor.CanPaste() },
			Do:    func(a any) { c.Editor.Paste() },
		},
		ContextMenuItem{
			Title: "Cut",
			Can:   func(any) bool { return c.Editor.CanCut() },
			Do:    func(a any) { c.Editor.Cut() },
		},
		ContextMenuItem{
			Title: "Delete",
			Can:   func(any) bool { return c.Editor.CanDelete() },
			Do:    func(a any) { c.Editor.Delete() },
		},
		ContextMenuItem{
			Title: "SelectAll",
			Can:   func(any) bool { return c.Editor.CanSelectAll() },
			Do:    func(a any) { c.Editor.SelectAll() },
		},
		ContextMenuItem{
			Title: "Save",
			// Can: func(any) bool{return multiLineField.CanCopy()},
			Do: func(a any) {
				return
				// os.Remove("")
			},
		},
		ContextMenuItem{
			Title: "SaveAs",
			Can:   nil,
			Do: func(a any) {
				return
				// os.Remove("")
			},
		},
		ContextMenuItem{
			Title: "Duplicate",
			Can:   nil,
			Do: func(a any) {
				return
				// os.Remove("")
			},
		},
		ContextMenuItem{
			Title: "open dir",
			Can:   nil,
			Do: func(a any) {
				stream.RunCommand("explorer") // todo
			},
		},
	).Install()
	c.Editor.RequestFocus()
	parent.AsPanel().AddChild(c)
	return c
}
