// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package widget

import (
	"github.com/ddkwork/unison"
	"github.com/ddkwork/unison/enums/paintstyle"
	"slices"
	"sort"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath"
)

// DefaultTableHeaderTheme holds the default TableHeaderTheme values for TableHeaders. Modifying this data will not
// alter existing TableHeaders, but will alter any TableHeaders created in the future.
var DefaultTableHeaderTheme = TableHeaderTheme{
	BackgroundInk:        unison.ControlColor,
	InteriorDividerColor: unison.InteriorDividerColor,
	HeaderBorder:         unison.NewLineBorder(unison.InteriorDividerColor, 0, unison.Insets{Bottom: 1}, false),
}

// TableHeaderTheme holds theming data for a TableHeader.
type TableHeaderTheme struct {
	BackgroundInk        unison.Ink
	InteriorDividerColor unison.Ink
	HeaderBorder         unison.Border
}

// TableHeader provides a header for a Node.
type TableHeader[T any] struct {
	unison.Panel
	TableHeaderTheme
	table                *Node[T]
	ColumnHeaders        []TableColumnHeader[T]
	Less                 func(s1, s2 string) bool
	interactionColumn    int
	columnResizeStart    float32
	columnResizeBase     float32
	columnResizeOverhead float32
	inHeader             bool
}

// NewTableHeader creates a new TableHeader.
func NewTableHeader[T any](table *Node[T], columnHeaders ...TableColumnHeader[T]) *TableHeader[T] {
	h := &TableHeader[T]{
		TableHeaderTheme: DefaultTableHeaderTheme,
		table:            table,
		ColumnHeaders:    columnHeaders,
		Less:             func(s1, s2 string) bool { return txt.NaturalLess(s1, s2, true) },
	}
	h.Self = h
	h.SetSizer(h.DefaultSizes)
	h.SetBorder(h.TableHeaderTheme.HeaderBorder)
	h.DrawCallback = h.DefaultDraw
	h.UpdateCursorCallback = h.DefaultUpdateCursorCallback
	h.UpdateTooltipCallback = h.DefaultUpdateTooltipCallback
	h.MouseMoveCallback = h.DefaultMouseMove
	h.MouseDownCallback = h.DefaultMouseDown
	h.MouseDragCallback = h.DefaultMouseDrag
	h.MouseUpCallback = h.DefaultMouseUp
	h.table.header = h
	return h
}

// DefaultSizes provides the default sizing.
func (h *TableHeader[T]) DefaultSizes(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
	prefSize.Width = h.table.FrameRect().Size.Width
	prefSize.Height = h.heightForColumns()
	if border := h.Border(); border != nil {
		insets := border.Insets()
		prefSize.Height += insets.Height()
	}
	return unison.NewSize(16, prefSize.Height), prefSize, prefSize
}

// ColumnFrame returns the frame of the given column.
func (h *TableHeader[T]) ColumnFrame(col int) unison.Rect {
	if col < 0 || col >= len(h.table.Columns) {
		return unison.Rect{}
	}
	insets := h.combinedInsets()
	x := insets.Left
	for c := 0; c < col; c++ {
		x += h.table.Columns[c].Current
		if h.table.ShowColumnDivider {
			x++
		}
	}
	rect := unison.NewRect(x, insets.Top, h.table.Columns[col].Current, h.FrameRect().Height-insets.Height())
	rect.Inset(h.table.Padding)
	return rect
}

func (h *TableHeader[T]) heightForColumns() float32 {
	var height float32
	for i := range h.table.Columns {
		w := h.table.Columns[i].Current
		if w <= 0 {
			continue
		}
		w -= h.table.Padding.Left + h.table.Padding.Right
		if i < len(h.ColumnHeaders) {
			_, cpref, _ := h.ColumnHeaders[i].AsPanel().Sizes(unison.Size{Width: w})
			cpref.Height += h.table.Padding.Top + h.table.Padding.Bottom
			if height < cpref.Height {
				height = cpref.Height
			}
		}
	}
	return max(xmath.Ceil(height), h.table.MinimumRowHeight)
}

func (h *TableHeader[T]) combinedInsets() unison.Insets {
	var insets unison.Insets
	if border := h.Border(); border != nil {
		insets = border.Insets()
	}
	if border := h.table.Border(); border != nil {
		insets2 := border.Insets()
		if insets.Left < insets2.Left {
			insets.Left = insets2.Left
		}
		if insets.Right < insets2.Right {
			insets.Right = insets2.Right
		}
	}
	return insets
}

// DefaultDraw provides the default drawing.
func (h *TableHeader[T]) DefaultDraw(canvas *unison.Canvas, dirty unison.Rect) {
	canvas.DrawRect(dirty, h.BackgroundInk.Paint(canvas, dirty, paintstyle.Fill))

	var firstCol int
	insets := h.combinedInsets()
	x := insets.Left
	for i := range h.table.Columns {
		x1 := x + h.table.Columns[i].Current
		if h.table.ShowColumnDivider {
			x1++
		}
		if x1 >= dirty.X {
			break
		}
		x = x1
		firstCol = i + 1
	}

	if h.table.ShowColumnDivider {
		rect := dirty
		rect.X = x
		rect.Width = 1
		for c := firstCol; c < len(h.table.Columns)-1; c++ {
			rect.X += h.table.Columns[c].Current
			canvas.DrawRect(rect, h.InteriorDividerColor.Paint(canvas, rect, paintstyle.Fill))
			rect.X++
		}
	}

	rect := dirty
	rect.X = x
	rect.Y = insets.Top
	rect.Height = h.heightForColumns()
	lastX := dirty.Right()
	for c := firstCol; c < len(h.table.Columns) && rect.X < lastX; c++ {
		rect.Width = h.table.Columns[c].Current
		cellRect := rect
		cellRect.Inset(h.table.Padding)
		if c < len(h.ColumnHeaders) {
			cell := h.ColumnHeaders[c].AsPanel()
			h.installCell(cell, cellRect)
			canvas.Save()
			canvas.Translate(cellRect.X, cellRect.Y)
			cellRect.X = 0
			cellRect.Y = 0
			cell.Draw(canvas, cellRect)
			h.uninstallCell(cell)
			canvas.Restore()
		}
		rect.X += h.table.Columns[c].Current
		if h.table.ShowColumnDivider {
			rect.X++
		}
	}
}

func (h *TableHeader[T]) installCell(cell *unison.Panel, frame unison.Rect) {
	cell.SetFrameRect(frame)
	cell.ValidateLayout()
	h.AsPanel().AddChild(cell)
}

func (h *TableHeader[T]) uninstallCell(cell *unison.Panel) {
	cell.RemoveFromParent()
}

// DefaultUpdateCursorCallback provides the default cursor update handling.
func (h *TableHeader[T]) DefaultUpdateCursorCallback(where unison.Point) *unison.Cursor {
	if !h.table.PreventUserColumnResize {
		if over := h.table.OverColumnDivider(where.X); over != -1 {
			if h.table.Columns[over].Minimum <= 0 || h.table.Columns[over].Minimum < h.table.Columns[over].Maximum {
				return unison.ResizeHorizontalCursor()
			}
		}
	}
	if col := h.table.OverColumn(where.X); col != -1 {
		cell := h.ColumnHeaders[col].AsPanel()
		if cell.UpdateCursorCallback != nil {
			rect := h.ColumnFrame(col)
			h.installCell(cell, rect)
			where.Subtract(rect.Point)
			cursor := cell.UpdateCursorCallback(where)
			h.uninstallCell(cell)
			return cursor
		}
	}
	return nil
}

// DefaultUpdateTooltipCallback provides the default tooltip update handling.
func (h *TableHeader[T]) DefaultUpdateTooltipCallback(where unison.Point, suggestedAvoidInRoot unison.Rect) unison.Rect {
	if col := h.table.OverColumn(where.X); col != -1 {
		cell := h.ColumnHeaders[col].AsPanel()
		if cell.UpdateTooltipCallback != nil {
			rect := h.ColumnFrame(col)
			h.installCell(cell, rect)
			where.Subtract(rect.Point)
			rect = h.RectToRoot(rect)
			rect.Align()
			avoid := cell.UpdateTooltipCallback(where, rect)
			h.Tooltip = cell.Tooltip
			h.uninstallCell(cell)
			return avoid
		}
		if cell.Tooltip != nil {
			h.Tooltip = cell.Tooltip
			suggestedAvoidInRoot = h.RectToRoot(h.ColumnFrame(col))
			suggestedAvoidInRoot.Align()
			return suggestedAvoidInRoot
		}
	}
	h.Tooltip = nil
	return unison.Rect{}
}

// DefaultMouseMove provides the default mouse move handling.
func (h *TableHeader[T]) DefaultMouseMove(where unison.Point, mod unison.Modifiers) bool {
	stop := false
	if col := h.table.OverColumn(where.X); col != -1 {
		cell := h.ColumnHeaders[col].AsPanel()
		if cell.MouseMoveCallback != nil {
			rect := h.ColumnFrame(col)
			h.installCell(cell, rect)
			where.Subtract(rect.Point)
			stop = cell.MouseMoveCallback(where, mod)
			h.uninstallCell(cell)
		}
	}
	return stop
}

// DefaultMouseDown provides the default mouse down handling.
func (h *TableHeader[T]) DefaultMouseDown(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
	h.interactionColumn = -1
	h.inHeader = false
	if !h.table.PreventUserColumnResize {
		if over := h.table.OverColumnDivider(where.X); over != -1 {
			if h.table.Columns[over].Minimum <= 0 || h.table.Columns[over].Minimum < h.table.Columns[over].Maximum {
				if clickCount == 2 {
					h.table.SizeColumnToFit(over, true)
					h.MarkForRedraw()
					h.Window().UpdateCursorNow()
					return true
				}
				h.interactionColumn = over
				h.columnResizeStart = where.X
				h.columnResizeBase = h.table.Columns[over].Current
				h.columnResizeOverhead = h.table.Padding.Left + h.table.Padding.Right
				if h.table.Columns[over].ID == h.table.HierarchyColumnID {
					depth := 0
					for _, cache := range h.table.rowCache {
						if depth < cache.depth {
							depth = cache.depth
						}
					}
					h.columnResizeOverhead += h.table.Padding.Left + h.table.HierarchyIndent*float32(depth+1)
				}
				return true
			}
		}
	}
	stop := true
	if col := h.table.OverColumn(where.X); col != -1 {
		h.interactionColumn = col
		h.inHeader = true
		cell := h.ColumnHeaders[col].AsPanel()
		if cell.MouseDownCallback != nil {
			rect := h.ColumnFrame(col)
			h.installCell(cell, rect)
			where.Subtract(rect.Point)
			stop = cell.MouseDownCallback(where, button, clickCount, mod)
			h.uninstallCell(cell)
		}
	}
	return stop
}

// DefaultMouseDrag provides the default mouse drag handling.
func (h *TableHeader[T]) DefaultMouseDrag(where unison.Point, _ int, _ unison.Modifiers) bool {
	if !h.table.PreventUserColumnResize && !h.inHeader && h.interactionColumn != -1 {
		width := h.columnResizeBase + where.X - h.columnResizeStart
		if width < h.columnResizeOverhead {
			width = h.columnResizeOverhead
		}
		minimum := h.table.Columns[h.interactionColumn].Minimum
		if minimum > 0 && width < minimum+h.columnResizeOverhead {
			width = minimum + h.columnResizeOverhead
		} else {
			maximum := h.table.Columns[h.interactionColumn].Maximum
			if maximum > 0 && width > maximum+h.columnResizeOverhead {
				width = maximum + h.columnResizeOverhead
			}
		}
		if h.table.Columns[h.interactionColumn].Current != width {
			h.table.Columns[h.interactionColumn].Current = width
			h.table.SyncToModel()
			h.MarkForRedraw()
		}
		return true
	}
	return false
}

// DefaultMouseUp provides the default mouse up handling.
func (h *TableHeader[T]) DefaultMouseUp(where unison.Point, button int, mod unison.Modifiers) bool {
	stop := false
	if h.inHeader && h.interactionColumn != -1 {
		cell := h.ColumnHeaders[h.interactionColumn].AsPanel()
		if cell.MouseUpCallback != nil {
			rect := h.ColumnFrame(h.interactionColumn)
			h.installCell(cell, rect)
			where.Subtract(rect.Point)
			stop = cell.MouseUpCallback(where, button, mod)
			h.uninstallCell(cell)
		}
	}
	return stop
}

// SortOn adjusts the sort such that the specified header is the primary sort column. If the header was already the
// primary sort column, then its ascending/descending flag will be flipped instead.
func (h *TableHeader[T]) SortOn(header TableColumnHeader[T]) {
	if header.SortState().Sortable {
		headers := make([]TableColumnHeader[T], len(h.ColumnHeaders))
		copy(headers, h.ColumnHeaders)
		sort.Slice(headers, func(i, j int) bool {
			if headers[i] == header {
				return true
			}
			if headers[j] == header {
				return false
			}
			s1 := headers[i].SortState()
			if !s1.Sortable || s1.Order < 0 {
				return false
			}
			s2 := headers[j].SortState()
			if !s2.Sortable || s2.Order < 0 {
				return true
			}
			return s1.Order < s2.Order
		})
		for i, hdr := range headers {
			s := hdr.SortState()
			if s.Sortable {
				if i == 0 {
					if s.Order == 0 {
						s.Ascending = !s.Ascending
					} else {
						s.Order = 0
					}
				} else if s.Order >= 0 {
					s.Order = i
				}
			} else {
				s.Order = -1
			}
			hdr.SetSortState(s)
		}
	}
}

type headerWithIndex[T any] struct {
	index  int
	header TableColumnHeader[T]
}

// HasSort returns true if at least one column is marked for sorting.
func (h *TableHeader[T]) HasSort() bool {
	for _, hdr := range h.ColumnHeaders {
		if ss := hdr.SortState(); ss.Sortable && ss.Order >= 0 {
			return true
		}
	}
	return false
}

// ApplySort sorts the table according to the current sort criteria.
func (h *TableHeader[T]) ApplySort() {
	headers := make([]*headerWithIndex[T], len(h.ColumnHeaders))
	for i, hdr := range h.ColumnHeaders {
		headers[i] = &headerWithIndex[T]{
			index:  i,
			header: hdr,
		}
	}
	sort.Slice(headers, func(i, j int) bool {
		s1 := headers[i].header.SortState()
		if !s1.Sortable || s1.Order < 0 {
			return false
		}
		s2 := headers[j].header.SortState()
		if !s2.Sortable || s2.Order < 0 {
			return true
		}
		return s1.Order < s2.Order
	})
	for i, hdr := range headers {
		s := hdr.header.SortState()
		if !s.Sortable || s.Order < 0 {
			headers = headers[:i]
			break
		}
	}
	if h.table.filteredRows == nil {
		roots := slices.Clone(h.table.RootRows())
		h.applySort(headers, roots)
		h.table.SetRootRows(roots) // Avoid resetting the selection by directly updating the model
	} else {
		h.applySort(headers, h.table.filteredRows)
	}
	h.table.SyncToModel()
}

func (h *TableHeader[T]) applySort(headers []*headerWithIndex[T], rows []*Node[T]) {
	if len(headers) > 0 && len(rows) > 0 {
		sort.Slice(rows, func(i, j int) bool {
			for _, hdr := range headers {
				d1 := rows[i].CellDataForSort(hdr.index)
				d2 := rows[j].CellDataForSort(hdr.index)
				if d1 != d2 {
					ascending := hdr.header.SortState().Ascending
					if h.Less(d1, d2) {
						return ascending
					}
					return !ascending
				}
			}
			return false
		})
		if h.table.filteredRows == nil {
			for _, row := range rows {
				if row.CanHaveChildren() {
					if children := row.Children; len(children) > 1 {
						children = slices.Clone(children)
						h.applySort(headers, children)
						row.SetChildren(children)
					}
				}
			}
		}
	}
}
