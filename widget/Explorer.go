package widget

import (
	"os"
	"path/filepath"
	"time"

	"github.com/richardwilkes/toolbox/i18n"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/datasize"
	"github.com/richardwilkes/unison"
)

// https://github.com/y4v8/filewatcher
// https://github.com/Ronbb/usn
type (
	Explorer struct {
		*unison.Dock
		Root *Node[DirTree]
	}
	DirTree struct {
		Name    string
		Size    int64
		Type    string
		ModTime time.Time
		Path    string
	}
)

func NewExplorer(parent unison.Paneler, walkDir string) unison.Paneler {
	e := &Explorer{
		Dock: unison.NewDock(),
		Root: nil,
	}
	e.Self = e

	table, header := NewTable(
		DirTree{
			Name:    "",
			Size:    0,
			Type:    "",
			ModTime: time.Time{},
			Path:    "",
		},
		TableContext[DirTree]{
			ContextMenuItems: nil,
			MarshalRow: func(node *Node[DirTree]) (cells []CellData) {
				if node.Container() {
					node.Data.Name = node.Data.Type
					sum := int64(0)
					node.Data.Size = sum
					node.Walk(func(node *Node[DirTree]) {
						sum += node.Data.Size
					})
					node.Data.Size = sum
				}
				return []CellData{
					{Text: node.Data.Name},
					{Text: datasize.Size(node.Data.Size).String()},
					{Text: node.Data.Type},
					{Text: stream.FormatTime(node.Data.ModTime)},
					{Text: node.Data.Path},
				}
			},
			UnmarshalRow:             nil,
			SelectionChangedCallback: nil,
			SetRootRowsCallBack: func(root *Node[DirTree]) {
				// todo 新建文本文档,dark title bar
				if walkDir == "" {
					walkDir = "."
				}
				e.Walk(walkDir, root)
				stream.WriteTruncate("explorer.txt", root.String())
				root.SetRootRows(root.Children)
				root.SizeColumnsToFit(true)
			},
			JsonName:   "explorer",
			IsDocument: false,
		},
	)
	return NewTableScrollPanel(parent, table, header)
}

func (e *Explorer) Walk(path string, parent *Node[DirTree]) (ok bool) {
	parent.Data.Type = filepath.Base(path) // 设置root的type，在格式化回调中赋值到Name
	dirEntries := mylog.Check2(os.ReadDir(path))
	for _, entry := range dirEntries {
		info := mylog.Check2(entry.Info())
		dirTree := DirTree{
			Name:    entry.Name(),
			Size:    info.Size(),
			Type:    entry.Type().String(),
			ModTime: info.ModTime(),
			Path:    filepath.Join(path, entry.Name()),
		}
		if entry.IsDir() {
			containerNode := NewContainerNode(entry.Name(), dirTree)
			parent.AddChild(containerNode)
			e.Walk(filepath.Join(path, entry.Name()), containerNode)
			continue
		}
		parent.AddChild(NewNode(dirTree))
	}
	return true
}

func (e *Explorer) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  unison.DocumentSVG,
		Size: suggestedSize,
	}
}

func (e *Explorer) Title() string {
	return i18n.Text("Document Workspace")
}

func (e *Explorer) Tooltip() string {
	return ""
}

func (e *Explorer) Modified() bool {
	return false
}
