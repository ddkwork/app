package widget

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/app"

	"github.com/ddkwork/golibrary/stream"
	"github.com/richardwilkes/unison/enums/align"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/google/uuid"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

var zeroUUID = uuid.UUID{}

type TableDragData[T any] struct {
	Table *Node[T]
	Rows  []*Node[T]
}

type ColumnInfo struct {
	ID          int
	Current     float32
	Minimum     float32
	Maximum     float32
	AutoMinimum float32
	AutoMaximum float32
	Cell        string
}

type tableCache[T any] struct {
	row    *Node[T]
	parent int
	depth  int
	height float32
}

type tableHitRect struct {
	unison.Rect
	handler func()
}

var DefaultTableTheme = TableTheme{
	BackgroundInk:          unison.ThemeBelowSurface,
	OnBackgroundInk:        unison.ThemeOnBelowSurface,
	BandingInk:             unison.ThemeSurface,
	OnBandingInk:           unison.ThemeOnSurface,
	InteriorDividerInk:     unison.ThemeAboveSurface,
	SelectionInk:           unison.ThemeFocus,
	OnSelectionInk:         unison.ThemeOnFocus,
	InactiveSelectionInk:   unison.ThemeDeepFocus,
	OnInactiveSelectionInk: unison.ThemeOnDeepFocus,
	IndirectSelectionInk:   unison.ThemeDeeperFocus,
	OnIndirectSelectionInk: unison.ThemeOnDeeperFocus,
	Padding:                unison.NewUniformInsets(4),
	HierarchyIndent:        16,
	MinimumRowHeight:       16,
	ColumnResizeSlop:       4,
	ShowColumnDivider:      true,
}

type TableTheme struct {
	BackgroundInk          unison.Ink    `json:"-"`
	OnBackgroundInk        unison.Ink    `json:"-"`
	BandingInk             unison.Ink    `json:"-"`
	OnBandingInk           unison.Ink    `json:"-"`
	InteriorDividerInk     unison.Ink    `json:"-"`
	SelectionInk           unison.Ink    `json:"-"`
	OnSelectionInk         unison.Ink    `json:"-"`
	InactiveSelectionInk   unison.Ink    `json:"-"`
	OnInactiveSelectionInk unison.Ink    `json:"-"`
	IndirectSelectionInk   unison.Ink    `json:"-"`
	OnIndirectSelectionInk unison.Ink    `json:"-"`
	Padding                unison.Insets `json:"-"`
	HierarchyColumnID      int           `json:"-"`
	HierarchyIndent        float32       `json:"-"`
	MinimumRowHeight       float32       `json:"-"`
	ColumnResizeSlop       float32       `json:"-"`
	ShowRowDivider         bool          `json:"-"`
	ShowColumnDivider      bool          `json:"-"`
}
type CellData struct {
	Text     string
	MaxWidth float32
	Disabled bool
	Tooltip  string

	SvgBuffer   string
	ImageBuffer []byte
	FgColor     unison.Color
}

const ContainerKeyPostfix = "_container"

type TableContext[T any] struct {
	ContextMenuItems         func(node *Node[T]) []ContextMenuItem
	MarshalRow               func(node *Node[T]) (cells []CellData)
	UnmarshalRow             func(node *Node[T], values []string)
	SelectionChangedCallback func(root *Node[T])
	SetRootRowsCallBack      func(root *Node[T])
	JsonName                 string
	IsDocument               bool
}
type Node[T any] struct {
	unison.Panel `json:"-"`
	TableTheme   `json:"-"`
	MarshalRow   func(node *Node[T]) (cells []CellData) `json:"-"`

	ID         uuid.UUID      `json:"id"`
	Type       string         `json:"type"`
	ThirdParty map[string]any `json:"-"`
	parent     *Node[T]
	Data       T
	Children   []*Node[T] `json:"children,omitempty"`
	isOpen     bool       `json:"open,omitempty"`

	SelectionChangedCallback func() `json:"-"`
	DoubleClickCallback      func() `json:"-"`
	DragRemovedRowsCallback  func() `json:"-"`
	DropOccurredCallback     func() `json:"-"`
	Columns                  []ColumnInfo
	filteredRows             []*Node[T]
	header                   *TableHeader[T]
	selMap                   map[uuid.UUID]bool
	selAnchor                uuid.UUID
	lastSel                  uuid.UUID
	hitRects                 []tableHitRect
	rowCache                 []tableCache[T]
	lastMouseEnterCellPanel  *unison.Panel
	lastMouseDownCellPanel   *unison.Panel
	interactionRow           int
	interactionColumn        int
	lastMouseMotionRow       int
	lastMouseMotionColumn    int
	startRow                 int
	endBeforeRow             int
	columnResizeStart        float32
	columnResizeBase         float32
	columnResizeOverhead     float32
	PreventUserColumnResize  bool
	awaitingSizeColumnsToFit bool
	awaitingSyncToModel      bool
	selNeedsPrune            bool
	wasDragged               bool
	dividerDrag              bool
}

func NewTableScrollPanel[T any](table *Node[T], header *TableHeader[T]) *unison.Panel {
	panel := NewPanel()
	panel.AddChild(table)
	panel.AddChild(header)
	scrollPanelFill := NewScrollPanelFill(panel)
	scrollPanelFill.SetColumnHeader(header)
	return scrollPanelFill.AsPanel()
}

func NewTableScroll[T any](data T, ctx TableContext[T]) *unison.Panel {
	table, header := NewTable[T](data, ctx)
	content := NewPanel()
	content.AddChild(table)
	content.AddChild(header)
	scrollPanelFill := NewScrollPanelFill(content)
	scrollPanelFill.SetColumnHeader(header)
	return scrollPanelFill.AsPanel()
}

func (n *Node[T]) AddChildByData(data T) { n.AddChild(NewNode(data)) }
func (n *Node[T]) AddChildByDatas(datas ...T) {
	for _, data := range datas {
		n.AddChild(NewNode(data))
	}
}

func (n *Node[T]) AddContainerByData(typeKey string, data T) (newContainer *Node[T]) { // 我们需要返回新的容器节点用于递归填充它的孩子节点，用例是explorer文件资源管理器
	newContainer = NewContainerNode(typeKey, data)
	n.AddChild(newContainer)
	return newContainer
}

func (n *Node[T]) Sum(key string) string {
	// container column 0 key is empty string
	key = n.Type
	key = strings.TrimSuffix(key, ContainerKeyPostfix)
	key += " (" + fmt.Sprint(n.LenChildren()) + ")"
	return key
}

func (n *Node[T]) CopyColumn() string {
	b := stream.NewBuffer("var columnData = []string{")
	b.NewLine()
	n.Walk(func(node *Node[T]) {
		cells := n.MarshalRow(node)
		b.WriteString(strconv.Quote(cells[n.header.interactionColumn].Text))
		b.WriteStringLn(",")
	})
	b.WriteStringLn("}")
	unison.GlobalClipboard.SetText(b.String())
	return b.String()
}

func (n *Node[T]) CopyRow() string {
	b := stream.NewBuffer("var rowData = []string{")
	rows := n.SelectedRows(false)
	for _, row := range rows {
		cells := row.MarshalRow(row)
		for i, cell := range cells {
			b.WriteString(strconv.Quote(cell.Text))
			if i < len(cells)-1 {
				b.WriteString(",")
			}
		}
	}
	b.WriteStringLn("}")
	unison.GlobalClipboard.SetText(b.String())
	return b.String()
}

func (n *Node[T]) CloneForTarget(target unison.Paneler, newParent *Node[T]) *Node[T] {
	mylog.Todo("remove target")
	clone := *n
	clone.parent = newParent
	clone.ID = uuid.New()
	return &clone
}

func NewUUID() uuid.UUID {
	return mylog.Check2(uuid.NewRandom())
}

func (n *Node[T]) UUID() uuid.UUID {
	return n.ID
}

func (n *Node[T]) Container() bool {
	return strings.HasSuffix(n.Type, ContainerKeyPostfix)
}

func (n *Node[T]) kind(base string) string {
	if n.Container() {
		return fmt.Sprintf(i18n.Text("%s Container"), base)
	}
	return base
}

func (n *Node[T]) GetType() string {
	return n.Type
}

func (n *Node[T]) SetType(typeKey string) {
	n.Type = typeKey
}

func (n *Node[T]) IsOpen() bool {
	return n.isOpen && n.Container()
}

func (n *Node[T]) SetOpen(open bool) {
	n.isOpen = open && n.Container()
}

func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

func (n *Node[T]) SetParent(parent *Node[T]) {
	n.parent = parent
}

func (n *Node[T]) clearUnusedFields() {
	if !n.Container() {
		n.Children = nil
		n.isOpen = false
	}
}

func (n *Node[T]) CanHaveChildren() bool {
	return n.HasChildren()
}

func (n *Node[T]) HasChildren() bool {
	return n.Container() && len(n.Children) > 0
}

func (n *Node[T]) SetChildren(children []*Node[T]) {
	n.Children = children
}

func (n *Node[T]) CellDataForSort(col int) string {
	return n.MarshalRow(n)[col].Text
}

func (n *Node[T]) ColumnCell(row, col int, foreground, background unison.Ink, selected, indirectlySelected, focused bool) unison.Paneler {
	wrapper := unison.NewPanel()
	cells := n.MarshalRow(n)
	wrapper.SetLayout(&unison.FlexLayout{Columns: 1})
	width := n.CellWidth(row, col)

	maxWidth := float32(30)
	if width > maxWidth {
		width = maxWidth
	} else {
		maxWidth = width
	}
	data := CellData{
		Text:        cells[col].Text,
		MaxWidth:    maxWidth,
		Disabled:    cells[col].Disabled,
		Tooltip:     cells[col].Tooltip,
		SvgBuffer:   cells[col].SvgBuffer,
		ImageBuffer: cells[col].ImageBuffer,
		FgColor:     cells[col].FgColor,
	}
	addWrappedText(wrapper, foreground, unison.LabelFont, data)

	wrapper.UpdateTooltipCallback = func(_ unison.Point, _ unison.Rect) unison.Rect {
		wrapper.Tooltip = unison.NewTooltipWithText("A tooltip for the cell")
		return wrapper.RectToRoot(wrapper.ContentRect(true))
	}
	return wrapper
}

func addWrappedText(parent *unison.Panel, ink unison.Ink, font unison.Font, data CellData) {
	decoration := &unison.TextDecoration{Font: font}
	var lines []*unison.Text
	if data.MaxWidth > 0 {
		lines = unison.NewTextWrappedLines(data.Text, decoration, data.MaxWidth)
	} else {
		lines = unison.NewTextLines(data.Text, decoration)
	}
	for _, line := range lines {
		label := unison.NewLabel()
		label.SetTitle(line.String())
		label.Font = font
		label.LabelTheme.OnBackgroundInk = ink
		if data.Disabled {
			label.LabelTheme.BackgroundInk = unison.DarkRed
		}
		size := unison.LabelFont.Size() + 7

		if data.ImageBuffer != nil {
			label.Drawable = &unison.SizedDrawable{
				Drawable: nil,
				Size:     unison.NewSize(size, size),
			}
		}
		if data.SvgBuffer != "" {
			label.Drawable = &unison.DrawableSVG{
				SVG:  nil,
				Size: unison.NewSize(size, size),
			}
		}
		// LabelStyle(label)
		parent.AddChild(label)
	}
}

func initHeader(data any) (Columns []ColumnInfo) {
	fields := stream.ReflectVisibleFields(data)
	Columns = make([]ColumnInfo, 0, len(fields))
	for i, field := range fields {
		label := unison.NewLabel()
		label.SetTitle(field.Name)
		Columns = append(Columns, ColumnInfo{
			ID:          i,
			Current:     0,
			Minimum:     20,
			Maximum:     10000,
			AutoMinimum: 0,
			AutoMaximum: 0,
			Cell:        field.Name,
		})
	}
	return
}

func NewNode[T any](data T) (child *Node[T]) {
	return newNode("", false, data)
}

func NewContainerNode[T any](typeKey string, data T) (container *Node[T]) {
	n := newNode(typeKey, true, data)
	n.Children = make([]*Node[T], 0)
	return n
}

func NewContainerNodes[T any](typeKeys []string, objects ...T) (containerNodes []*Node[T]) {
	containerNodes = make([]*Node[T], 0)
	var data T // it is zero value
	for i, key := range typeKeys {
		if len(objects) > 0 {
			data = objects[i]
		}
		containerNodes = append(containerNodes, NewContainerNode(key, data))
	}
	return
}

func NewTable[T any](data T, ctx TableContext[T]) (table *Node[T], header *TableHeader[T]) {
	if ctx.JsonName == "" {
		mylog.Check("JsonName is empty")
	}
	ctx.JsonName = strings.TrimSuffix(ctx.JsonName, ".json")
	table, header = newTable(data, ctx)
	fnUpdate := func() {
		table.SetRootRows(table.Children)
		table.SizeColumnsToFit(true)
		stream.MarshalJsonToFile(table, ctx.JsonName+".json")
		stream.WriteTruncate(ctx.JsonName+".txt", table.Document())
		if ctx.IsDocument {
			b := stream.NewBuffer("")
			b.WriteStringLn("# " + ctx.JsonName + " document table")
			b.WriteStringLn("```text")
			b.WriteStringLn(table.Document())
			b.WriteStringLn("```")
			stream.WriteTruncate("README.md", b.String())
		}
	}
	switch {
	case ctx.SetRootRowsCallBack != nil:
		ctx.SetRootRowsCallBack(table)
		fallthrough
	case ctx.SelectionChangedCallback != nil:
		table.SelectionChangedCallback = func() { ctx.SelectionChangedCallback(table) }
		fallthrough
	case table.FileDropCallback == nil:
		table.FileDropCallback = func(files []string) {
			if filepath.Ext(files[0]) == ".json" {
				table.ResetChildren()
				b := stream.NewBuffer(files[0])
				mylog.Check(json.Unmarshal(b.Bytes(), table)) // todo test need a zero table?
				fnUpdate()
			}
			mylog.Struct(files)
		}
		fallthrough
	default:
		table.DoubleClickCallback = func() {
			rows := table.SelectedRows(false)
			for i, row := range rows {
				app.Run("edit row #"+fmt.Sprint(i), func(w *unison.Window) {
					content := w.Content()
					nodeEditor, RowPanel := NewStructView(row.Data, func(data T) (values []CellData) {
						return table.MarshalRow(row)
					})
					content.AddChild(nodeEditor)
					content.AddChild(RowPanel)
					panel := NewButtonsPanel(
						[]string{
							"apply", "cancel",
						},
						func() { //todo bug,need merge strctView here for apply
							mylog.CheckNil(ctx.UnmarshalRow)                   //noce edit must implement UnmarshalRow callback
							ctx.UnmarshalRow(row, nodeEditor.getFieldValues()) //replace nodeEditor UnmarshalRow
							nodeEditor.Update(row.Data)
							table.SyncToModel()
							stream.MarshalJsonToFile(table.Children, ctx.JsonName+".json")
							w.Dispose()
						},
						func() {
							w.Dispose()
						},
					)
					RowPanel.AddChild(panel)
					RowPanel.AddChild(NewVSpacer())
				})
			}
		}
	}
	fnUpdate()
	return
}

func newTable[T any](data T, ctx TableContext[T]) (table *Node[T], header *TableHeader[T]) {
	root := NewContainerNode("root", data)
	root.MarshalRow = ctx.MarshalRow
	// root.contextMenuItems = ctx.ContextMenuItems
	// root.root = root
	root.Columns = initHeader(data)

	for i, column := range root.Columns {
		text := unison.NewText(column.Cell, &unison.TextDecoration{
			Font:            unison.LabelFont,
			BackgroundInk:   nil,
			OnBackgroundInk: nil,
			BaselineOffset:  0,
			Underline:       false,
			StrikeThrough:   false,
		})
		root.Columns[i].Minimum = text.Width() + root.Padding.Left + root.Padding.Right
	}

	root.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
		if mod == 0 && (keyCode == unison.KeyBackspace || keyCode == unison.KeyDelete) {
			root.PerformCmd(root, unison.DeleteItemID)
			return true
		}
		return root.DefaultKeyDown(keyCode, mod, repeat) // todo add delete,move to ctx menu,exporter need delete file or dir
	}

	if ctx.ContextMenuItems != nil {
		contextMenuItems := ctx.ContextMenuItems(root)
		NewContextMenuItems(root, root.DefaultMouseDown, contextMenuItems...).Install()
	}

	root.InstallDragSupport(nil, "dragKey", "singularName", "pluralName")
	InstallDropSupport[T, any](root, "dragKey", func(from, to *Node[T]) bool { return from == to }, nil, nil)

	header = NewTableHeader[T](root)
	for _, column := range root.Columns {
		columnHeader := NewTableColumnHeader[T](column.Cell, "")
		columnHeader.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
			return true
		}
		NewContextMenuItems(columnHeader, columnHeader.Label.MouseDownCallback,
			ContextMenuItem{
				Title: "copy column",
				id:    0,
				Can: func(a any) bool {
					return true
				},
				Do: func(a any) { root.CopyColumn() },
			},
		).Install()
		header.ColumnHeaders = append(header.ColumnHeaders, columnHeader)
	}

	header.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
	})
	return root, header
}

func (n *Node[T]) IsRoot() bool { return n.parent == nil }
func newNode[T any](typeKey string, isContainer bool, data T) *Node[T] {
	if isContainer {
		typeKey += ContainerKeyPostfix
	}

	n := &Node[T]{
		Panel:                    unison.Panel{},
		TableTheme:               DefaultTableTheme,
		MarshalRow:               nil,
		ID:                       NewUUID(),
		Type:                     typeKey,
		ThirdParty:               nil,
		parent:                   nil,
		Data:                     data,
		Children:                 nil,
		isOpen:                   isContainer,
		SelectionChangedCallback: nil,
		DoubleClickCallback:      nil,
		DragRemovedRowsCallback:  nil,
		DropOccurredCallback:     nil,
		Columns:                  nil,
		filteredRows:             nil,
		header:                   nil,
		selMap:                   make(map[uuid.UUID]bool),
		selAnchor:                uuid.UUID{},
		lastSel:                  uuid.UUID{},
		hitRects:                 nil,
		rowCache:                 nil,
		lastMouseEnterCellPanel:  nil,
		lastMouseDownCellPanel:   nil,
		interactionRow:           -1,
		interactionColumn:        -1,
		lastMouseMotionRow:       -1,
		lastMouseMotionColumn:    -1,
		startRow:                 0,
		endBeforeRow:             0,
		columnResizeStart:        0,
		columnResizeBase:         0,
		columnResizeOverhead:     0,
		PreventUserColumnResize:  false,
		awaitingSizeColumnsToFit: false,
		awaitingSyncToModel:      false,
		selNeedsPrune:            false,
		wasDragged:               false,
		dividerDrag:              false,
	}

	n.Self = n
	n.SetFocusable(true)
	n.SetSizer(n.DefaultSizes)
	n.GainedFocusCallback = n.DefaultFocusGained
	n.DrawCallback = n.DefaultDraw
	n.UpdateCursorCallback = n.DefaultUpdateCursorCallback
	n.UpdateTooltipCallback = n.DefaultUpdateTooltipCallback
	n.MouseMoveCallback = n.DefaultMouseMove
	n.MouseDownCallback = n.DefaultMouseDown
	n.MouseDragCallback = n.DefaultMouseDrag
	n.MouseUpCallback = n.DefaultMouseUp
	n.MouseEnterCallback = n.DefaultMouseEnter
	n.MouseExitCallback = n.DefaultMouseExit
	n.KeyDownCallback = n.DefaultKeyDown
	n.InstallCmdHandlers(unison.SelectAllItemID, unison.AlwaysEnabled, func(_ any) { n.SelectAll() })
	InstallContainerConversionHandlers(n, n)
	n.wasDragged = false

	NewContextMenuItems(n, n.DefaultMouseDown,
		ContextMenuItem{
			Title: "CopyRow",
			id:    0,
			Can:   func(a any) bool { return true },
			Do:    func(a any) { n.CopyRow() },
		},
		ContextMenuItem{
			Title: "ConvertToContainer",
			id:    0,
			Can:   func(a any) bool { return CanConvertToContainer(n) },
			Do:    func(a any) { ConvertToContainer(n) },
		},
		ContextMenuItem{
			Title: "ConvertToNonContainer",
			id:    0,
			Can:   func(a any) bool { return CanConvertToNonContainer(n) },
			Do:    func(a any) { ConvertToNonContainer(n) },
		},
		ContextMenuItem{
			Title: "NewNode",
			id:    0,
			Can:   func(a any) bool { return true },
			Do: func(a any) {
				node := NewNode(data)
				rows := n.SelectedRows(false)
				for _, row := range rows {
					row.AddChild(node)
					n.SyncToModel()
				}
			},
		},
		ContextMenuItem{
			Title: "NewContainerNode",
			id:    0,
			Can:   func(a any) bool { return true },
			Do: func(a any) {
				container := NewContainerNode(n.Type, data)
				rows := n.SelectedRows(false)
				for _, row := range rows {
					row.AddChild(container)
					n.SyncToModel()
				}
			},
		},
		ContextMenuItem{
			Title: "DeleteNode",
			id:    0,
			Can:   func(a any) bool { return true },
			Do: func(a any) {
				rows := n.SelectedRows(false)
				for i, row := range rows {
					mylog.Todo("bug: deleting a node doesn'n update the model")
					slices.Delete(row.Children, i, i+1)
				}
				n.SyncToModel()
			},
		},
		ContextMenuItem{
			Title: "duplicateNode",
			id:    0,
			Can:   func(a any) bool { return true },
			Do: func(a any) {
				rows := n.SelectedRows(false)
				for _, row := range rows {
					row.AddChild(row.CloneForTarget(row.AsPanel(), row))
				}
				n.SyncToModel()
			},
		},
		ContextMenuItem{
			Title: "EditNode",
			id:    0,
			Can:   func(a any) bool { return true },
			Do:    func(a any) { mylog.Todo("implement edit node") },
		},
	).Install()

	return n
}

func (n *Node[T]) AddChild(child *Node[T]) {
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *Node[T]) ColumnIndexForID(id int) int {
	for i, c := range n.Columns {
		if c.ID == id {
			return i
		}
	}
	return -1
}

func (n *Node[T]) SetDrawRowRange(start, endBefore int) {
	n.startRow = start
	n.endBeforeRow = endBefore
}

func (n *Node[T]) ClearDrawRowRange() {
	n.startRow = 0
	n.endBeforeRow = 0
}

func (n *Node[T]) CurrentDrawRowRange() (start, endBefore int) {
	if n.startRow < n.endBeforeRow && n.startRow >= 0 && n.endBeforeRow <= len(n.rowCache) {
		return n.startRow, n.endBeforeRow
	}
	return 0, len(n.rowCache)
}

func (n *Node[T]) DefaultDraw(canvas *unison.Canvas, dirty unison.Rect) {
	selectionInk := n.SelectionInk
	if !n.Focused() {
		selectionInk = n.InactiveSelectionInk
	}

	canvas.DrawRect(dirty, n.BackgroundInk.Paint(canvas, dirty, paintstyle.Fill))

	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}

	var firstCol int
	x := insets.Left
	for i := range n.Columns {
		x1 := x + n.Columns[i].Current
		if n.ShowColumnDivider {
			x1++
		}
		if x1 >= dirty.X {
			break
		}
		x = x1
		firstCol = i + 1
	}

	startRow, endBeforeRow := n.CurrentDrawRowRange()
	y := insets.Top
	for i := startRow; i < endBeforeRow; i++ {
		y1 := y + n.rowCache[i].height
		if n.ShowRowDivider {
			y1++
		}
		if y1 >= dirty.Y {
			break
		}
		y = y1
		startRow = i + 1
	}

	lastY := dirty.Bottom()
	rect := dirty
	rect.Y = y
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ {
		rect.Height = n.rowCache[r].height
		if n.IsRowOrAnyParentSelected(r) {
			if n.IsRowSelected(r) {
				canvas.DrawRect(rect, selectionInk.Paint(canvas, rect, paintstyle.Fill))
			} else {
				canvas.DrawRect(rect, n.IndirectSelectionInk.Paint(canvas, rect, paintstyle.Fill))
			}
		} else if r%2 == 1 {
			canvas.DrawRect(rect, n.BandingInk.Paint(canvas, rect, paintstyle.Fill))
		}
		rect.Y += n.rowCache[r].height
		if n.ShowRowDivider && r != endBeforeRow-1 {
			rect.Height = 1
			canvas.DrawRect(rect, n.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill))
			rect.Y++
		}
	}

	if n.ShowColumnDivider {
		rect = dirty
		rect.X = x
		rect.Width = 1
		for c := firstCol; c < len(n.Columns)-1; c++ {
			rect.X += n.Columns[c].Current
			canvas.DrawRect(rect, n.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill))
			rect.X++
		}
	}

	rect = dirty
	rect.Y = y
	lastX := dirty.Right()
	n.hitRects = nil
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ {
		rect.X = x
		rect.Height = n.rowCache[r].height
		for c := firstCol; c < len(n.Columns) && rect.X < lastX; c++ {
			fg, bg, selected, indirectlySelected, focused := n.cellParams(r, c)
			rect.Width = n.Columns[c].Current
			cellRect := rect.Inset(n.Padding)
			row := n.rowCache[r].row
			if n.Columns[c].ID == n.HierarchyColumnID {
				if row.CanHaveChildren() {
					const disclosureIndent = 2
					disclosureSize := min(n.HierarchyIndent, n.MinimumRowHeight) - disclosureIndent*2
					canvas.Save()
					left := cellRect.X + n.HierarchyIndent*float32(n.rowCache[r].depth) + disclosureIndent
					top := cellRect.Y + (n.MinimumRowHeight-disclosureSize)/2
					dSize := unison.Size{Width: disclosureSize, Height: disclosureSize}
					n.hitRects = append(n.hitRects,
						n.newTableHitRect(unison.Rect{Point: unison.Point{X: left, Y: top}, Size: dSize}, row))
					canvas.Translate(left, top)
					if row.IsOpen() {
						offset := disclosureSize / 2
						canvas.Translate(offset, offset)
						canvas.Rotate(90)
						canvas.Translate(-offset, -offset)
					}
					canvas.DrawPath(unison.CircledChevronRightSVG.PathForSize(dSize),
						fg.Paint(canvas, cellRect, paintstyle.Fill))
					canvas.Restore()
				}
				indent := n.HierarchyIndent*float32(n.rowCache[r].depth+1) + n.Padding.Left
				cellRect.X += indent
				cellRect.Width -= indent
			}
			cell := row.ColumnCell(r, c, fg, bg, selected, indirectlySelected, focused).AsPanel()
			n.installCell(cell, cellRect)
			canvas.Save()
			canvas.Translate(cellRect.X, cellRect.Y)
			cellRect.X = 0
			cellRect.Y = 0
			cell.Draw(canvas, cellRect)
			n.uninstallCell(cell)
			canvas.Restore()
			rect.X += n.Columns[c].Current
			if n.ShowColumnDivider {
				rect.X++
			}
		}
		rect.Y += n.rowCache[r].height
		if n.ShowRowDivider {
			rect.Y++
		}
	}
}

func (n *Node[T]) cellParams(row, _ int) (fg, bg unison.Ink, selected, indirectlySelected, focused bool) {
	focused = n.Focused()
	selected = n.IsRowSelected(row)
	indirectlySelected = !selected && n.IsRowOrAnyParentSelected(row)
	switch {
	case selected && focused:
		fg = n.OnSelectionInk
		bg = n.SelectionInk
	case selected:
		fg = n.OnInactiveSelectionInk
		bg = n.InactiveSelectionInk
	case indirectlySelected:
		fg = n.OnIndirectSelectionInk
		bg = n.IndirectSelectionInk
	case row%2 == 1:
		fg = n.OnBandingInk
		bg = n.BandingInk
	default:
		fg = n.OnBackgroundInk
		bg = n.BackgroundInk
	}
	return fg, bg, selected, indirectlySelected, focused
}

func (n *Node[T]) cell(row, col int) *unison.Panel {
	fg, bg, selected, indirectlySelected, focused := n.cellParams(row, col)
	return n.rowCache[row].row.ColumnCell(row, col, fg, bg, selected, indirectlySelected, focused).AsPanel()
}

func (n *Node[T]) installCell(cell *unison.Panel, frame unison.Rect) {
	cell.SetFrameRect(frame)
	cell.ValidateLayout()
	n.AsPanel().AddChild(cell)
}

func (n *Node[T]) uninstallCell(cell *unison.Panel) {
	cell.RemoveFromParent()
}

func (n *Node[T]) RowHeights() []float32 {
	heights := make([]float32, len(n.rowCache))
	for i := range n.rowCache {
		heights[i] = n.rowCache[i].height
	}
	return heights
}

func (n *Node[T]) OverRow(y float32) int {
	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}
	end := insets.Top
	for i := range n.rowCache {
		start := end
		end += n.rowCache[i].height
		if n.ShowRowDivider {
			end++
		}
		if y >= start && y < end {
			return i
		}
	}
	return -1
}

func (n *Node[T]) OverColumn(x float32) int {
	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}
	end := insets.Left
	for i := range n.Columns {
		start := end
		end += n.Columns[i].Current
		if n.ShowColumnDivider {
			end++
		}
		if x >= start && x < end {
			return i
		}
	}
	return -1
}

func (n *Node[T]) OverColumnDivider(x float32) int {
	if len(n.Columns) < 2 {
		return -1
	}
	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}
	pos := insets.Left
	for i := range n.Columns[:len(n.Columns)-1] {
		pos += n.Columns[i].Current
		if n.ShowColumnDivider {
			pos++
		}
		if xmath.Abs(pos-x) < n.ColumnResizeSlop {
			return i
		}
	}
	return -1
}

func (n *Node[T]) CellWidth(row, col int) float32 {
	if row < 0 || col < 0 || row >= len(n.rowCache) || col >= len(n.Columns) {
		return 0
	}
	width := n.Columns[col].Current - (n.Padding.Left + n.Padding.Right)
	if n.Columns[col].ID == n.HierarchyColumnID {
		width -= n.HierarchyIndent*float32(n.rowCache[row].depth+1) + n.Padding.Left
	}
	return width
}

func (n *Node[T]) ColumnEdges(col int) (left, right float32) {
	if col < 0 || col >= len(n.Columns) {
		return 0, 0
	}
	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}
	left = insets.Left
	for c := 0; c < col; c++ {
		left += n.Columns[c].Current
		if n.ShowColumnDivider {
			left++
		}
	}
	right = left + n.Columns[col].Current
	left += n.Padding.Left
	right -= n.Padding.Right
	if n.Columns[col].ID == n.HierarchyColumnID {
		left += n.HierarchyIndent + n.Padding.Left
	}
	if right < left {
		right = left
	}
	return left, right
}

func (n *Node[T]) CellFrame(row, col int) unison.Rect {
	if row < 0 || col < 0 || row >= len(n.rowCache) || col >= len(n.Columns) {
		return unison.Rect{}
	}
	var insets unison.Insets
	if border := n.Border(); border != nil {
		insets = border.Insets()
	}
	x := insets.Left
	for c := 0; c < col; c++ {
		x += n.Columns[c].Current
		if n.ShowColumnDivider {
			x++
		}
	}
	y := insets.Top
	for r := 0; r < row; r++ {
		y += n.rowCache[r].height
		if n.ShowRowDivider {
			y++
		}
	}
	rect := unison.Rect{
		Point: unison.Point{X: x, Y: y},
		Size:  unison.Size{Width: n.Columns[col].Current, Height: n.rowCache[row].height},
	}.Inset(n.Padding)
	if n.Columns[col].ID == n.HierarchyColumnID {
		indent := n.HierarchyIndent*float32(n.rowCache[row].depth+1) + n.Padding.Left
		rect.X += indent
		rect.Width -= indent
		if rect.Width < 1 {
			rect.Width = 1
		}
	}
	return rect
}

func (n *Node[T]) RowFrame(row int) unison.Rect {
	if row < 0 || row >= len(n.rowCache) {
		return unison.Rect{}
	}
	rect := n.ContentRect(false)
	for i := 0; i < row; i++ {
		rect.Y += n.rowCache[i].height
		if n.ShowRowDivider {
			rect.Y++
		}
	}
	rect.Height = n.rowCache[row].height
	return rect
}

func (n *Node[T]) newTableHitRect(rect unison.Rect, row *Node[T]) tableHitRect {
	return tableHitRect{
		Rect: rect,
		handler: func() {
			open := !row.IsOpen()
			row.SetOpen(open)
			n.SyncToModel()
			if !open {
				n.PruneSelectionOfUndisclosedNodes()
			}
		},
	}
}

func (n *Node[T]) DefaultFocusGained() {
	switch {
	case n.interactionRow != -1:
		n.ScrollRowIntoView(n.interactionRow)
	case n.lastMouseMotionRow != -1:
		n.ScrollRowIntoView(n.lastMouseMotionRow)
	default:
		n.ScrollIntoView()
	}
	n.MarkForRedraw()
}

func (n *Node[T]) DefaultUpdateCursorCallback(where unison.Point) *unison.Cursor {
	if !n.PreventUserColumnResize {
		if over := n.OverColumnDivider(where.X); over != -1 {
			if n.Columns[over].Minimum <= 0 || n.Columns[over].Minimum < n.Columns[over].Maximum {
				return unison.ResizeHorizontalCursor()
			}
		}
	}
	if row := n.OverRow(where.Y); row != -1 {
		if col := n.OverColumn(where.X); col != -1 {
			cell := n.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *unison.Panel) bool { return p.UpdateCursorCallback != nil }) {
				var cursor *unison.Cursor
				rect := n.CellFrame(row, col)
				n.installCell(cell, rect)
				where = where.Sub(rect.Point)
				target := cell.PanelAt(where)
				for target != n.AsPanel() {
					if target.UpdateCursorCallback == nil {
						target = target.Parent()
					} else {
						toolbox.Call(func() { cursor = target.UpdateCursorCallback(cell.PointTo(where, target)) })
						break
					}
				}
				n.uninstallCell(cell)
				return cursor
			}
		}
	}
	return nil
}

func (n *Node[T]) DefaultUpdateTooltipCallback(where unison.Point, avoid unison.Rect) unison.Rect {
	if row := n.OverRow(where.Y); row != -1 {
		if col := n.OverColumn(where.X); col != -1 {
			cell := n.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *unison.Panel) bool { return p.UpdateTooltipCallback != nil || p.Tooltip != nil }) {
				rect := n.CellFrame(row, col)
				n.installCell(cell, rect)
				where = where.Sub(rect.Point)
				target := cell.PanelAt(where)
				n.Tooltip = nil
				n.TooltipImmediate = false
				for target != n.AsPanel() {
					avoid = target.RectToRoot(target.ContentRect(true)).Align()
					if target.UpdateTooltipCallback != nil {
						toolbox.Call(func() { avoid = target.UpdateTooltipCallback(cell.PointTo(where, target), avoid) })
					}
					if target.Tooltip != nil {
						n.Tooltip = target.Tooltip
						n.TooltipImmediate = target.TooltipImmediate
						break
					}

					target = target.Parent()
				}
				n.uninstallCell(cell)
				return avoid
			}
			if cell.Tooltip != nil {
				n.Tooltip = cell.Tooltip
				n.TooltipImmediate = cell.TooltipImmediate
				return n.RectToRoot(n.CellFrame(row, col)).Align()
			}
		}
	}
	n.Tooltip = nil
	return unison.Rect{}
}

func (n *Node[T]) DefaultMouseEnter(where unison.Point, mod unison.Modifiers) bool {
	row := n.OverRow(where.Y)
	col := n.OverColumn(where.X)
	if n.lastMouseMotionRow != row || n.lastMouseMotionColumn != col {
		n.DefaultMouseExit()
		n.lastMouseMotionRow = row
		n.lastMouseMotionColumn = col
	}
	if row != -1 && col != -1 {
		cell := n.cell(row, col)
		rect := n.CellFrame(row, col)
		n.installCell(cell, rect)
		where = where.Sub(rect.Point)
		target := cell.PanelAt(where)
		if target != n.lastMouseEnterCellPanel && n.lastMouseEnterCellPanel != nil {
			n.DefaultMouseExit()
			n.lastMouseMotionRow = row
			n.lastMouseMotionColumn = col
		}
		if target.MouseEnterCallback != nil {
			toolbox.Call(func() { target.MouseEnterCallback(cell.PointTo(where, target), mod) })
		}
		n.uninstallCell(cell)
		n.lastMouseEnterCellPanel = target
	}
	return true
}

func (n *Node[T]) DefaultMouseMove(where unison.Point, mod unison.Modifiers) bool {
	n.DefaultMouseEnter(where, mod)
	if n.lastMouseEnterCellPanel != nil {
		row := n.OverRow(where.Y)
		col := n.OverColumn(where.X)
		cell := n.cell(row, col)
		rect := n.CellFrame(row, col)
		n.installCell(cell, rect)
		where = where.Sub(rect.Point)
		if target := cell.PanelAt(where); target.MouseMoveCallback != nil {
			toolbox.Call(func() { target.MouseMoveCallback(cell.PointTo(where, target), mod) })
		}
		n.uninstallCell(cell)
	}
	return true
}

func (n *Node[T]) DefaultMouseExit() bool {
	if n.lastMouseEnterCellPanel != nil && n.lastMouseEnterCellPanel.MouseExitCallback != nil &&
		n.lastMouseMotionColumn != -1 && n.lastMouseMotionRow >= 0 && n.lastMouseMotionRow < len(n.rowCache) {
		cell := n.cell(n.lastMouseMotionRow, n.lastMouseMotionColumn)
		rect := n.CellFrame(n.lastMouseMotionRow, n.lastMouseMotionColumn)
		n.installCell(cell, rect)
		toolbox.Call(func() { n.lastMouseEnterCellPanel.MouseExitCallback() })
		n.uninstallCell(cell)
	}
	n.lastMouseEnterCellPanel = nil
	n.lastMouseMotionRow = -1
	n.lastMouseMotionColumn = -1
	return true
}

func (n *Node[T]) DefaultMouseDown(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
	if n.Window().InDrag() {
		return false
	}
	n.RequestFocus()
	n.wasDragged = false
	n.dividerDrag = false
	n.lastSel = zeroUUID

	n.interactionRow = -1
	n.interactionColumn = -1
	if button == unison.ButtonLeft {
		if !n.PreventUserColumnResize {
			if over := n.OverColumnDivider(where.X); over != -1 {
				if n.Columns[over].Minimum <= 0 || n.Columns[over].Minimum < n.Columns[over].Maximum {
					if clickCount == 2 {
						n.SizeColumnToFit(over, true)
						n.MarkForRedraw()
						n.Window().UpdateCursorNow()
						return true
					}
					n.interactionColumn = over
					n.columnResizeStart = where.X
					n.columnResizeBase = n.Columns[over].Current
					n.columnResizeOverhead = n.Padding.Left + n.Padding.Right
					if n.Columns[over].ID == n.HierarchyColumnID {
						depth := 0
						for _, cache := range n.rowCache {
							if depth < cache.depth {
								depth = cache.depth
							}
						}
						n.columnResizeOverhead += n.Padding.Left + n.HierarchyIndent*float32(depth+1)
					}
					return true
				}
			}
		}
		for _, one := range n.hitRects {
			if where.In(one.Rect) {
				return true
			}
		}
	}
	if row := n.OverRow(where.Y); row != -1 {
		if col := n.OverColumn(where.X); col != -1 {
			cell := n.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *unison.Panel) bool { return p.MouseDownCallback != nil }) {
				n.interactionRow = row
				n.interactionColumn = col
				rect := n.CellFrame(row, col)
				n.installCell(cell, rect)
				where = where.Sub(rect.Point)
				stop := false
				if target := cell.PanelAt(where); target.MouseDownCallback != nil {
					n.lastMouseDownCellPanel = target
					toolbox.Call(func() {
						stop = target.MouseDownCallback(cell.PointTo(where, target), button,
							clickCount, mod)
					})
				}
				n.uninstallCell(cell)
				if stop {
					return stop
				}
			}
		}
		rowData := n.rowCache[row].row
		id := rowData.UUID()
		switch {
		case mod&unison.ShiftModifier != 0:
			selAnchorIndex := -1
			if n.selAnchor != zeroUUID {
				for i, c := range n.rowCache {
					if c.row.UUID() == n.selAnchor {
						selAnchorIndex = i
						break
					}
				}
			}
			if selAnchorIndex != -1 {
				last := max(selAnchorIndex, row)
				for i := min(selAnchorIndex, row); i <= last; i++ {
					n.selMap[n.rowCache[i].row.UUID()] = true
				}
				n.notifyOfSelectionChange()
			} else if !n.selMap[id] {
				n.selMap = make(map[uuid.UUID]bool)
				n.selMap[id] = true
				n.selAnchor = id
				n.notifyOfSelectionChange()
			}
		case mod.DiscontiguousSelectionDown():
			if n.selMap[id] {
				delete(n.selMap, id)
			} else {
				n.selMap[id] = true
			}
			n.notifyOfSelectionChange()
		case n.selMap[id]:
			n.lastSel = id
		default:
			n.selMap = make(map[uuid.UUID]bool)
			n.selMap[id] = true
			n.selAnchor = id
			n.notifyOfSelectionChange()
		}
		n.MarkForRedraw()
		if button == unison.ButtonLeft && clickCount == 2 && n.DoubleClickCallback != nil && len(n.selMap) != 0 {
			toolbox.Call(n.DoubleClickCallback)
		}
	}
	return true
}

func (n *Node[T]) notifyOfSelectionChange() {
	if n.SelectionChangedCallback != nil {
		toolbox.Call(n.SelectionChangedCallback)
	}
}

func (n *Node[T]) DefaultMouseDrag(where unison.Point, button int, mod unison.Modifiers) bool {
	n.wasDragged = true
	stop := false
	if n.interactionColumn != -1 {
		if n.interactionRow == -1 {
			if button == unison.ButtonLeft && !n.PreventUserColumnResize {
				width := n.columnResizeBase + where.X - n.columnResizeStart
				if width < n.columnResizeOverhead {
					width = n.columnResizeOverhead
				}
				minimum := n.Columns[n.interactionColumn].Minimum
				if minimum > 0 && width < minimum+n.columnResizeOverhead {
					width = minimum + n.columnResizeOverhead
				} else {
					maximum := n.Columns[n.interactionColumn].Maximum
					if maximum > 0 && width > maximum+n.columnResizeOverhead {
						width = maximum + n.columnResizeOverhead
					}
				}
				if n.Columns[n.interactionColumn].Current != width {
					n.Columns[n.interactionColumn].Current = width
					n.EventuallySyncToModel()
					n.MarkForRedraw()
					n.dividerDrag = true
				}
				stop = true
			}
		} else if n.lastMouseDownCellPanel != nil && n.lastMouseDownCellPanel.MouseDragCallback != nil {
			cell := n.cell(n.interactionRow, n.interactionColumn)
			rect := n.CellFrame(n.interactionRow, n.interactionColumn)
			n.installCell(cell, rect)
			where = where.Sub(rect.Point)
			toolbox.Call(func() {
				stop = n.lastMouseDownCellPanel.MouseDragCallback(cell.PointTo(where, n.lastMouseDownCellPanel), button, mod)
			})
			n.uninstallCell(cell)
		}
	}
	return stop
}

func (n *Node[T]) DefaultMouseUp(where unison.Point, button int, mod unison.Modifiers) bool {
	stop := false
	if !n.dividerDrag && button == unison.ButtonLeft {
		for _, one := range n.hitRects {
			if where.In(one.Rect) {
				one.handler()
				stop = true
				break
			}
		}
	}

	if !n.wasDragged && n.lastSel != zeroUUID {
		n.ClearSelection()
		n.selMap[n.lastSel] = true
		n.selAnchor = n.lastSel
		n.MarkForRedraw()
		n.notifyOfSelectionChange()
	}

	if !stop && n.interactionRow != -1 && n.interactionColumn != -1 && n.lastMouseDownCellPanel != nil &&
		n.lastMouseDownCellPanel.MouseUpCallback != nil {
		cell := n.cell(n.interactionRow, n.interactionColumn)
		rect := n.CellFrame(n.interactionRow, n.interactionColumn)
		n.installCell(cell, rect)
		where = where.Sub(rect.Point)
		toolbox.Call(func() {
			stop = n.lastMouseDownCellPanel.MouseUpCallback(cell.PointTo(where, n.lastMouseDownCellPanel), button, mod)
		})
		n.uninstallCell(cell)
	}
	n.lastMouseDownCellPanel = nil
	n.interactionRow = -1
	return stop
}

func (n *Node[T]) DefaultKeyDown(keyCode unison.KeyCode, mod unison.Modifiers, _ bool) bool {
	if unison.IsControlAction(keyCode, mod) {
		if n.DoubleClickCallback != nil && len(n.selMap) != 0 {
			toolbox.Call(n.DoubleClickCallback)
		}
		return true
	}
	switch keyCode {
	case unison.KeyLeft:
		if n.HasSelection() {
			altered := false
			for _, row := range n.SelectedRows(false) {
				if row.IsOpen() {
					row.SetOpen(false)
					altered = true
				}
			}
			if altered {
				n.SyncToModel()
				n.PruneSelectionOfUndisclosedNodes()
			}
		}
	case unison.KeyRight:
		if n.HasSelection() {
			altered := false
			for _, row := range n.SelectedRows(false) {
				if !row.IsOpen() {
					row.SetOpen(true)
					altered = true
				}
			}
			if altered {
				n.SyncToModel()
			}
		}
	case unison.KeyUp:
		var i int
		if n.HasSelection() {
			i = max(n.FirstSelectedRowIndex()-1, 0)
		} else {
			i = len(n.rowCache) - 1
		}
		if !mod.ShiftDown() {
			n.ClearSelection()
		}
		n.SelectByIndex(i)
		n.ScrollRowCellIntoView(i, 0)
	case unison.KeyDown:
		i := min(n.LastSelectedRowIndex()+1, len(n.rowCache)-1)
		if !mod.ShiftDown() {
			n.ClearSelection()
		}
		n.SelectByIndex(i)
		n.ScrollRowCellIntoView(i, 0)
	case unison.KeyHome:
		if mod.ShiftDown() && n.HasSelection() {
			n.SelectRange(0, n.FirstSelectedRowIndex())
		} else {
			n.ClearSelection()
			n.SelectByIndex(0)
		}
		n.ScrollRowCellIntoView(0, 0)
	case unison.KeyEnd:
		if mod.ShiftDown() && n.HasSelection() {
			n.SelectRange(n.LastSelectedRowIndex(), len(n.rowCache)-1)
		} else {
			n.ClearSelection()
			n.SelectByIndex(len(n.rowCache) - 1)
		}
		n.ScrollRowCellIntoView(len(n.rowCache)-1, 0)
	default:
		return false
	}
	return true
}

func (n *Node[T]) PruneSelectionOfUndisclosedNodes() {
	if !n.selNeedsPrune {
		return
	}
	n.selNeedsPrune = false
	if len(n.selMap) == 0 {
		return
	}
	needsNotify := false
	selMap := make(map[uuid.UUID]bool, len(n.selMap))
	for _, entry := range n.rowCache {
		id := entry.row.UUID()
		if n.selMap[id] {
			selMap[id] = true
		} else {
			needsNotify = true
		}
	}
	n.selMap = selMap
	if needsNotify {
		n.notifyOfSelectionChange()
	}
}

func (n *Node[T]) FirstSelectedRowIndex() int {
	if len(n.selMap) == 0 {
		return -1
	}
	for i, entry := range n.rowCache {
		if n.selMap[entry.row.UUID()] {
			return i
		}
	}
	return -1
}

func (n *Node[T]) LastSelectedRowIndex() int {
	if len(n.selMap) == 0 {
		return -1
	}
	for i := len(n.rowCache) - 1; i >= 0; i-- {
		if n.selMap[n.rowCache[i].row.UUID()] {
			return i
		}
	}
	return -1
}

func (n *Node[T]) IsRowOrAnyParentSelected(index int) bool {
	if index < 0 || index >= len(n.rowCache) {
		return false
	}
	for index >= 0 {
		if n.selMap[n.rowCache[index].row.UUID()] {
			return true
		}
		index = n.rowCache[index].parent
	}
	return false
}

func (n *Node[T]) IsRowSelected(index int) bool {
	if index < 0 || index >= len(n.rowCache) {
		return false
	}
	return n.selMap[n.rowCache[index].row.UUID()]
}

func (n *Node[T]) SelectedRows(minimal bool) []*Node[T] {
	n.PruneSelectionOfUndisclosedNodes()
	if len(n.selMap) == 0 {
		return nil
	}
	rows := make([]*Node[T], 0, len(n.selMap))
	for _, entry := range n.rowCache {
		if n.selMap[entry.row.UUID()] && (!minimal || entry.parent == -1 || !n.IsRowOrAnyParentSelected(entry.parent)) {
			rows = append(rows, entry.row)
		}
	}
	return rows
}

func (n *Node[T]) CopySelectionMap() map[uuid.UUID]bool {
	n.PruneSelectionOfUndisclosedNodes()
	return copySelMap(n.selMap)
}

func (n *Node[T]) SetSelectionMap(selMap map[uuid.UUID]bool) {
	n.selMap = copySelMap(selMap)
	n.selNeedsPrune = true
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func copySelMap(selMap map[uuid.UUID]bool) map[uuid.UUID]bool {
	result := make(map[uuid.UUID]bool, len(selMap))
	for k, v := range selMap {
		result[k] = v
	}
	return result
}

func (n *Node[T]) HasSelection() bool {
	n.PruneSelectionOfUndisclosedNodes()
	return len(n.selMap) != 0
}

func (n *Node[T]) SelectionCount() int {
	n.PruneSelectionOfUndisclosedNodes()
	return len(n.selMap)
}

func (n *Node[T]) ClearSelection() {
	if len(n.selMap) == 0 {
		return
	}
	n.selMap = make(map[uuid.UUID]bool)
	n.selNeedsPrune = false
	n.selAnchor = zeroUUID
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) SelectAll() {
	n.selMap = make(map[uuid.UUID]bool, len(n.rowCache))
	n.selNeedsPrune = false
	n.selAnchor = zeroUUID
	for _, cache := range n.rowCache {
		id := cache.row.UUID()
		n.selMap[id] = true
		if n.selAnchor == zeroUUID {
			n.selAnchor = id
		}
	}
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) SelectByIndex(indexes ...int) {
	for _, index := range indexes {
		if index >= 0 && index < len(n.rowCache) {
			id := n.rowCache[index].row.UUID()
			n.selMap[id] = true
			n.selNeedsPrune = true
			if n.selAnchor == zeroUUID {
				n.selAnchor = id
			}
		}
	}
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) SelectRange(start, end int) {
	start = max(start, 0)
	end = min(end, len(n.rowCache)-1)
	if start > end {
		return
	}
	for i := start; i <= end; i++ {
		id := n.rowCache[i].row.UUID()
		n.selMap[id] = true
		n.selNeedsPrune = true
		if n.selAnchor == zeroUUID {
			n.selAnchor = id
		}
	}
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) DeselectByIndex(indexes ...int) {
	for _, index := range indexes {
		if index >= 0 && index < len(n.rowCache) {
			delete(n.selMap, n.rowCache[index].row.UUID())
		}
	}
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) DeselectRange(start, end int) {
	start = max(start, 0)
	end = min(end, len(n.rowCache)-1)
	if start > end {
		return
	}
	for i := start; i <= end; i++ {
		delete(n.selMap, n.rowCache[i].row.UUID())
	}
	n.MarkForRedraw()
	n.notifyOfSelectionChange()
}

func (n *Node[T]) DiscloseRow(row *Node[T], delaySync bool) bool {
	modified := false
	p := row.Parent()
	var zero *Node[T]
	for p != zero {
		if !p.IsOpen() {
			p.SetOpen(true)
			modified = true
		}
		p = p.Parent()
	}
	if modified {
		if delaySync {
			n.EventuallySyncToModel()
		} else {
			n.SyncToModel()
		}
	}
	return modified
}

func (n *Node[T]) RootRowCount() int {
	if n.filteredRows != nil {
		return len(n.filteredRows)
	}
	return n.RootRowCount()
}

func (n *Node[T]) RootRows() []*Node[T] {
	if n.filteredRows != nil {
		return n.filteredRows
	}
	return n.Children
}

func (n *Node[T]) SetRootRows(rows []*Node[T]) {
	n.filteredRows = nil
	n.Children = rows
	n.selMap = make(map[uuid.UUID]bool)
	n.selNeedsPrune = false
	n.selAnchor = zeroUUID
	n.SyncToModel()
}

func (n *Node[T]) SyncToModel() {
	rowCount := 0
	roots := n.RootRows()
	if n.filteredRows != nil {
		rowCount = len(n.filteredRows)
	} else {
		for _, row := range roots {
			rowCount += n.countOpenRowChildrenRecursively(row)
		}
	}
	n.rowCache = make([]tableCache[T], rowCount)
	j := 0
	for _, row := range roots {
		j = n.buildRowCacheEntry(row, -1, j, 0)
	}
	n.selNeedsPrune = true
	_, pref, _ := n.DefaultSizes(unison.Size{})
	rect := n.FrameRect()
	rect.Size = pref
	n.SetFrameRect(rect)
	n.MarkForRedraw()
	n.MarkForLayoutRecursivelyUpward()
}

func (n *Node[T]) countOpenRowChildrenRecursively(row *Node[T]) int {
	count := 1
	if row.CanHaveChildren() && row.IsOpen() {
		for _, child := range row.Children {
			count += n.countOpenRowChildrenRecursively(child)
		}
	}
	return count
}

func (n *Node[T]) buildRowCacheEntry(row *Node[T], parentIndex, index, depth int) int {
	row.MarshalRow = n.MarshalRow
	n.rowCache[index].row = row
	n.rowCache[index].parent = parentIndex
	n.rowCache[index].depth = depth
	n.rowCache[index].height = n.heightForColumns(row, index, depth)
	parentIndex = index
	index++
	if n.filteredRows == nil && row.CanHaveChildren() && row.IsOpen() {
		for _, child := range row.Children {
			index = n.buildRowCacheEntry(child, parentIndex, index, depth+1)
		}
	}
	return index
}

func (n *Node[T]) heightForColumns(rowData *Node[T], row, depth int) float32 {
	var height float32
	for col := range n.Columns {
		w := n.Columns[col].Current
		if w <= 0 {
			continue
		}
		w -= n.Padding.Left + n.Padding.Right
		if n.Columns[col].ID == n.HierarchyColumnID {
			w -= n.Padding.Left + n.HierarchyIndent*float32(depth+1)
		}
		size := n.cellPrefSize(rowData, row, col, w)
		size.Height += n.Padding.Top + n.Padding.Bottom
		if height < size.Height {
			height = size.Height
		}
	}
	return max(xmath.Ceil(height), n.MinimumRowHeight)
}

func (n *Node[T]) cellPrefSize(rowData *Node[T], row, col int, widthConstraint float32) unison.Size {
	fg, bg, selected, indirectlySelected, focused := n.cellParams(row, col)
	cell := rowData.ColumnCell(row, col, fg, bg, selected, indirectlySelected, focused).AsPanel()
	_, size, _ := cell.Sizes(unison.Size{Width: widthConstraint})
	return size
}

func (n *Node[T]) SizeColumnsToFitWithExcessIn(columnID int) {
	excessColumnIndex := max(n.ColumnIndexForID(columnID), 0)
	current := make([]float32, len(n.Columns))
	for col := range n.Columns {
		current[col] = max(n.Columns[col].Minimum, 0)
		n.Columns[col].Current = 0
	}
	for row, cache := range n.rowCache {
		for col := range n.Columns {
			if col == excessColumnIndex {
				continue
			}
			pref := n.cellPrefSize(cache.row, row, col, 0)
			minimum := n.Columns[col].AutoMinimum
			if minimum > 0 && pref.Width < minimum {
				pref.Width = minimum
			} else {
				maximum := n.Columns[col].AutoMaximum
				if maximum > 0 && pref.Width > maximum {
					pref.Width = maximum
				}
			}
			pref.Width += n.Padding.Left + n.Padding.Right
			if n.Columns[col].ID == n.HierarchyColumnID {
				pref.Width += n.Padding.Left + n.HierarchyIndent*float32(cache.depth+1)
			}
			if current[col] < pref.Width {
				current[col] = pref.Width
			}
		}
	}
	width := n.ContentRect(false).Width
	if n.ShowColumnDivider {
		width -= float32(len(n.Columns) - 1)
	}
	for col := range current {
		if col == excessColumnIndex {
			continue
		}
		n.Columns[col].Current = current[col]
		width -= current[col]
	}
	n.Columns[excessColumnIndex].Current = max(width, n.Columns[excessColumnIndex].Minimum)
	for row, cache := range n.rowCache {
		n.rowCache[row].height = n.heightForColumns(cache.row, row, cache.depth)
	}
}

func (n *Node[T]) SizeColumnsToFit(adjust bool) {
	current := make([]float32, len(n.Columns))
	for col := range n.Columns {
		current[col] = max(n.Columns[col].Minimum, 0)
		n.Columns[col].Current = 0
	}
	for row, cache := range n.rowCache {
		for col := range n.Columns {
			pref := n.cellPrefSize(cache.row, row, col, 0)
			minimum := n.Columns[col].AutoMinimum
			if minimum > 0 && pref.Width < minimum {
				pref.Width = minimum
			} else {
				maximum := n.Columns[col].AutoMaximum
				if maximum > 0 && pref.Width > maximum {
					pref.Width = maximum
				}
			}
			pref.Width += n.Padding.Left + n.Padding.Right
			if n.Columns[col].ID == n.HierarchyColumnID {
				pref.Width += n.Padding.Left + n.HierarchyIndent*float32(cache.depth+1)
			}
			if current[col] < pref.Width {
				current[col] = pref.Width
			}
		}
	}
	for col := range current {
		n.Columns[col].Current = current[col]
	}
	for row, cache := range n.rowCache {
		n.rowCache[row].height = n.heightForColumns(cache.row, row, cache.depth)
	}
	if adjust {
		_, pref, _ := n.DefaultSizes(unison.Size{})
		rect := n.FrameRect()
		rect.Size = pref
		n.SetFrameRect(rect)
	}
}

func (n *Node[T]) SizeColumnToFit(col int, adjust bool) {
	if col < 0 || col >= len(n.Columns) {
		return
	}
	current := max(n.Columns[col].Minimum, 0)
	n.Columns[col].Current = 0
	for row, cache := range n.rowCache {
		pref := n.cellPrefSize(cache.row, row, col, 0)
		minimum := n.Columns[col].AutoMinimum
		if minimum > 0 && pref.Width < minimum {
			pref.Width = minimum
		} else {
			maximum := n.Columns[col].AutoMaximum
			if maximum > 0 && pref.Width > maximum {
				pref.Width = maximum
			}
		}
		pref.Width += n.Padding.Left + n.Padding.Right
		if n.Columns[col].ID == n.HierarchyColumnID {
			pref.Width += n.Padding.Left + n.HierarchyIndent*float32(cache.depth+1)
		}
		if current < pref.Width {
			current = pref.Width
		}
	}
	n.Columns[col].Current = current
	for row, cache := range n.rowCache {
		n.rowCache[row].height = n.heightForColumns(cache.row, row, cache.depth)
	}
	if adjust {
		_, pref, _ := n.DefaultSizes(unison.Size{})
		rect := n.FrameRect()
		rect.Size = pref
		n.SetFrameRect(rect)
	}
}

func (n *Node[T]) EventuallySizeColumnsToFit(adjust bool) {
	if !n.awaitingSizeColumnsToFit {
		n.awaitingSizeColumnsToFit = true
		unison.InvokeTaskAfter(func() {
			n.SizeColumnsToFit(adjust)
			n.awaitingSizeColumnsToFit = false
		}, 20*time.Millisecond)
	}
}

func (n *Node[T]) EventuallySyncToModel() {
	if !n.awaitingSyncToModel {
		n.awaitingSyncToModel = true
		unison.InvokeTaskAfter(func() {
			n.SyncToModel()
			n.awaitingSyncToModel = false
		}, 20*time.Millisecond)
	}
}

func (n *Node[T]) DefaultSizes(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
	for col := range n.Columns {
		prefSize.Width += n.Columns[col].Current
	}
	startRow, endBeforeRow := n.CurrentDrawRowRange()
	for _, cache := range n.rowCache[startRow:endBeforeRow] {
		prefSize.Height += cache.height
	}
	if n.ShowColumnDivider {
		prefSize.Width += float32(len(n.Columns) - 1)
	}
	if n.ShowRowDivider {
		prefSize.Height += float32((endBeforeRow - startRow) - 1)
	}
	if border := n.Border(); border != nil {
		prefSize = prefSize.Add(border.Insets().Size())
	}
	prefSize = prefSize.Ceil()
	return prefSize, prefSize, prefSize
}

func (n *Node[T]) RowFromIndex(index int) *Node[T] {
	if index < 0 || index >= len(n.rowCache) {
		var zero *Node[T]
		return zero
	}
	return n.rowCache[index].row
}

func (n *Node[T]) RowToIndex(rowData *Node[T]) int {
	id := rowData.UUID()
	for row, data := range n.rowCache {
		if data.row.UUID() == id {
			return row
		}
	}
	return -1
}

func (n *Node[T]) LastRowIndex() int {
	return len(n.rowCache) - 1
}

func (n *Node[T]) ScrollRowIntoView(row int) {
	if frame := n.RowFrame(row); !frame.Empty() {
		n.ScrollRectIntoView(frame)
	}
}

func (n *Node[T]) ScrollRowCellIntoView(row, col int) {
	if frame := n.CellFrame(row, col); !frame.Empty() {
		n.ScrollRectIntoView(frame)
	}
}

func (n *Node[T]) IsFiltered() bool {
	return n.filteredRows != nil
}

func (n *Node[T]) ApplyFilter(filter func(row *Node[T]) bool) {
	if filter == nil {
		if n.filteredRows == nil {
			return
		}
		n.filteredRows = nil
	} else {
		n.filteredRows = make([]*Node[T], 0)
		for _, row := range n.RootRows() {
			n.applyFilter(row, filter)
		}
	}
	n.SyncToModel()
	if n.header != nil && n.header.HasSort() {
		n.header.ApplySort()
	}
}

func (n *Node[T]) applyFilter(row *Node[T], filter func(row *Node[T]) bool) {
	if !filter(row) {
		n.filteredRows = append(n.filteredRows, row)
	}
	if row.CanHaveChildren() {
		for _, child := range row.Children {
			n.applyFilter(child, filter)
		}
	}
}

func (n *Node[T]) InstallDragSupport(svg *unison.SVG, dragKey, singularName, pluralName string) {
	orig := n.MouseDragCallback
	n.MouseDragCallback = func(where unison.Point, button int, mod unison.Modifiers) bool {
		if orig != nil && orig(where, button, mod) {
			return true
		}
		if button == unison.ButtonLeft && n.HasSelection() && n.IsDragGesture(where) {
			data := &TableDragData[T]{
				Table: n,
				Rows:  n.SelectedRows(true),
			}
			drawable := NewTableDragDrawable(data, svg, singularName, pluralName)
			size := drawable.LogicalSize()
			n.StartDataDrag(&unison.DragData{
				Data:     map[string]any{dragKey: data},
				Drawable: drawable,
				Ink:      n.OnBackgroundInk,
				Offset:   unison.Point{X: 0, Y: -size.Height / 2},
			})
		}
		return false
	}
}

func InstallDropSupport[T any, U any](t *Node[T], dragKey string, shouldMoveDataCallback func(from, to *Node[T]) bool, willDropCallback func(from, to *Node[T], move bool) *unison.UndoEdit[U], didDropCallback func(undo *unison.UndoEdit[U], from, to *Node[T], move bool)) *TableDrop[T, U] {
	drop := &TableDrop[T, U]{
		Table:                  t,
		DragKey:                dragKey,
		originalDrawOver:       t.DrawOverCallback,
		shouldMoveDataCallback: shouldMoveDataCallback,
		willDropCallback:       willDropCallback,
		didDropCallback:        didDropCallback,
	}
	t.DataDragOverCallback = drop.DataDragOverCallback
	t.DataDragExitCallback = drop.DataDragExitCallback
	t.DataDragDropCallback = drop.DataDragDropCallback
	t.DrawOverCallback = drop.DrawOverCallback
	return drop
}

func CountTableRows[T any](rows []*Node[T]) int {
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children)
		}
	}
	return count
}

func RowContainsRow[T any](ancestor, descendant *Node[T]) bool {
	var zero *Node[T]
	for descendant != zero && descendant != ancestor {
		descendant = descendant.Parent()
	}
	return descendant == ancestor
}

func (n *Node[T]) RemoveFromParent() {
	mylog.CheckNil(n.parent)
	n.parent.Remove(n.ID)
}

func (n *Node[T]) Remove(id uuid.UUID) {
	for i, child := range n.Children {
		if child.ID == id {
			n.Children = slices.Delete(n.Children, i, i+1)
			break
		}
	}
}

func (n *Node[T]) Find(id uuid.UUID) *Node[T] {
	if n.ID == id {
		return n
	}
	for _, child := range n.Children {
		found := child.Find(id)
		if found != nil {
			return found
		}
	}
	return nil
}

func (n *Node[T]) Sort(cmp func(a T, b T) bool) {
	sort.SliceStable(n.Children, func(i, j int) bool {
		return cmp(n.Children[i].Data, n.Children[j].Data)
	})
	for _, child := range n.Children {
		child.Sort(cmp)
	}
}

func (n *Node[T]) Walk(callback func(node *Node[T])) {
	callback(n)
	for _, child := range n.Children {
		child.Walk(callback)
	}
}

func (n *Node[T]) WalkQueue(callback func(node *Node[T])) {
	queue := []*Node[T]{n}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		callback(node)
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
}

func (n *Node[T]) WalkContainer(callback func(node *Node[T])) {
	callback(n) // always walk root here
	containers := make([]*Node[T], 0)
	for _, child := range n.Children {
		if child.Container() {
			containers = append(containers, child)
		}
	}
	for _, container := range containers {
		container.Walk(callback)
	}
}

func (n *Node[T]) ApplyFilter_(tag string) {
	n.filteredRows = make([]*Node[T], 0)
	// var node *Node[T]
	// node = n.Root()

	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			cells := n.MarshalRow(node)
			for _, cell := range cells {
				if strings.EqualFold(cell.Text, tag) {
					n.filteredRows = append(n.filteredRows, node) // 先过滤所有容器节点
				}
			}
		}
	})

	for i, row := range n.filteredRows {
		children := make([]*Node[T], 0)
		row.Walk(func(node *Node[T]) {
			cells := n.MarshalRow(node)
			for _, cell := range cells {
				if strings.EqualFold(cell.Text, tag) {
					children = append(children, node) // 过滤子节点
				}
			}
		})
		n.filteredRows[i].SetChildren(children)
	}

	n.SetChildren(n.filteredRows)
}

type (
	TableTui struct {
		Header  Row      // 下面的行列操作来刷新表头的每列宽度，第一次不用取最大宽度，后面每刷新一行都要更新列宽
		Columns []Column // 动态添加行的同时填充每列的单元格切片并刷新最大宽度
		Rows    []Row    // 每添加一行，每个单元格的宽度来自列切片计算后的最大宽度，同时复制给表头的列宽
	}
	Row struct {
		Widths []int // every cell len,每行个每个单元格的宽度将在同时填充列的时候取通过所有行构造的列的最大宽度覆盖它
		Cells  []string
	}
	Column struct {
		MaxWidth int      // max(len(cells))
		Cells    []string // 每添加一行填充对应类的单元格数据并刷新最大宽度
	}
	/*
			header Column1 Column2 Column3 Column4
			row1   cell1   cell2   cell3   cell4    每插入一行，列切片的单元格切片增加一个，宽度取罗每列的单元格切片的宽度切片的最大值
			row2   cell1   cell2   cell3   cell4    每插入一行，每个单元格的宽度都来自每列的单元格切片的宽度的最大值
			...

			结构定义
			Columns()
			header()

			addRow() 同时填充每列的单元格切片的宽度和数据，并比较当前列的宽度，取最大的覆盖当前列的单元格宽度

			总的调整列宽算法就是：
		       每行的每个单元格宽度来自每列的所有单元格的最大宽度，
		       行列数据结构需要像上面定义的那样才方便理解。
	*/
)

func (n *Node[T]) Format(node *Node[T], s *stream.Buffer, isTui bool) TableTui {
	fields := stream.ReflectVisibleFields(n.Data)
	size := len(fields)
	tui := TableTui{
		Header: Row{
			Widths: make([]int, size),
			Cells:  make([]string, size),
		},
		Columns: make([]Column, size),
		Rows:    make([]Row, 0),
	}
	for i, field := range fields {
		tui.Header.Cells[i] = field.Name
		tui.Header.Widths[i] = len(field.Name)
		tui.Columns[i].Cells = append(tui.Columns[i].Cells, field.Name)
		tui.Columns[i].MaxWidth = len(field.Name)
	}

	// 1 渲染层级，递归格式化树节点
	const (
		indent          = "│   "
		childPrefix     = "├───"
		lastChildPrefix = "└───"
	)
	gioIndent := "    "

	tui.Rows = append(tui.Rows, tui.Header)

	node.Walk(func(node *Node[T]) { // 从根节点开始遍历
		Hierarchical := ""

		depth := node.Depth() - 1 // 根节点深度为0，每一层向下递增
		for i := 0; i < depth; i++ {
			// s.WriteString(indent) // 添加缩进
			if isTui {
				Hierarchical += indent // todo 为什么太宽？
			} else {
				Hierarchical += gioIndent
			}
		}

		if node.IsRoot() { // 添加节点前缀
			Hierarchical = "│" + Hierarchical
		} else if node.parent != nil && !node.IsLastChild() {
			if isTui {
				Hierarchical += childPrefix
			} else {
				Hierarchical += gioIndent
			}
		} else if node.parent != nil && node.IsLastChild() {
			if isTui {
				Hierarchical += lastChildPrefix
			} else {
				Hierarchical += gioIndent
			}
		}

		// 2 渲染行，添加节点数据
		// s.WriteString(fmt.Sprintf("[%v] Type: %v, Open: %v, Data: %v\n", node.UUID, node.Type, node.IsOpen, node.Data))
		mylog.CheckNil(n.MarshalRow)
		cells := n.MarshalRow(node) // 获取每行的单元格数据
		if len(cells) == 0 {        // 快速测试模式，业务模型还没建立好，树形还没准备好久跑单元测试的情况
			return
		}
		cells[0].Text = Hierarchical + cells[0].Text
		Hierarchical = ""
		row := Row{
			Widths: make([]int, len(cells)),
			Cells:  make([]string, len(cells)),
		}
		for i, cell := range cells {
			row.Cells[i] = cell.Text
			row.Widths[i] = tui.Columns[i].MaxWidth
			tui.Columns[i].Cells = append(tui.Columns[i].Cells, cell.Text)
			if len(cell.Text) > tui.Columns[i].MaxWidth {
				row.Widths[i] = len(cell.Text)
				tui.Columns[i].MaxWidth = len(cell.Text)
			}
		}
		tui.Rows = append(tui.Rows, row)
	})

	for index, row := range tui.Rows {
		if index == 0 {
			fnFmtHeader := func() (h string) {
				for i, cell := range row.Cells {
					if i < len(n.Columns)-1 {
						if i == 0 {
							indentStr := fmt.Sprintf("│%-*s ", tui.Columns[i].MaxWidth-1, cell) // 为什么要-1？
							h += indentStr
						} else {
							indentStr := fmt.Sprintf("│%-*s ", tui.Columns[i].MaxWidth, cell)
							h += indentStr
						}
					}
				}
				return
			}
			fmtHeader := fnFmtHeader()
			s.WriteStringLn(strings.Repeat("─", len(fmtHeader))) // 这个表头的矩形有点糟糕
			s.WriteStringLn(fmtHeader)
			s.WriteString(strings.Repeat("─", len(fmtHeader))) // todo 这里为什么不能换行？
			s.NewLine()
			continue
		}
		for i, cell := range row.Cells {
			if i < len(n.Columns)-1 {
				indentStr := fmt.Sprintf("%-*s ", tui.Columns[i].MaxWidth, cell) // 层级列已经有了层级文本了，不需要填充
				if i > 0 {
					indentStr = fmt.Sprintf("│%-*s ", tui.Columns[i].MaxWidth, cell) // 层级列之外需要列分隔符
				}
				s.WriteString(indentStr)
			}
		}
		s.NewLine()
	}
	return tui
}

func (n *Node[T]) String() string {
	s := stream.NewBuffer("")
	n.Format(n, s, true)
	return s.String()
}

func (n *Node[T]) Document() string {
	s := stream.NewBuffer("")
	// s.WriteStringLn("// interface or method name here")
	// s.WriteStringLn("/*")
	lines := stream.NewBuffer(n.String()).ToLines()
	for _, line := range lines {
		s.WriteStringLn("  " + line)
	}
	// s.WriteStringLn("*/")
	return s.String()
}

func (n *Node[T]) Depth() int {
	count := 0
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}

func (n *Node[T]) LenChildren() int {
	return len(n.Children)
}

func (n *Node[T]) LastChild() (lastChild *Node[T]) {
	if n.IsRoot() {
		return n.Children[len(n.Children)-1]
	}
	return n.parent.Children[len(n.parent.Children)-1]
}

func (n *Node[T]) IsLastChild() bool {
	return n.LastChild() == n
}

func (n *Node[T]) ResetChildren() {
	n.Children = nil
	n.rowCache = nil
	n.filteredRows = nil
}

func (n *Node[T]) OpenAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(true)
		}
	})
}

func (n *Node[T]) CloseAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(false)
		}
	})
}

func (n *Node[T]) CopyFrom(from *Node[T]) { // todo remove
	*n = *from
}

func (n *Node[T]) ApplyTo(to *Node[T]) { // todo remove
	*to = *n
}

func (n *Node[T]) Clone() (newNode *Node[T]) {
	defer n.SyncToModel()
	if n.Container() {
		return NewContainerNode(n.Type, n.Data)
	}
	return NewNode(n.Data)
}
