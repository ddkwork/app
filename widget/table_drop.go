package widget

import (
	"slices"

	"github.com/richardwilkes/unison"

	"github.com/google/uuid"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

type TableDrop[T any, U any] struct {
	Table                  *Node[T]
	DragKey                string
	TargetParent           *Node[T]
	TargetIndex            int
	AllDragData            map[string]any
	TableDragData          *TableDragData[T]
	originalDrawOver       func(*unison.Canvas, unison.Rect)
	shouldMoveDataCallback func(from, to *Node[T]) bool
	willDropCallback       func(from, to *Node[T], move bool) *unison.UndoEdit[U]
	didDropCallback        func(undo *unison.UndoEdit[U], from, to *Node[T], move bool)
	top                    float32
	left                   float32
	inDragOver             bool
}

func (d *TableDrop[T, U]) DrawOverCallback(gc *unison.Canvas, rect unison.Rect) {
	if d.originalDrawOver != nil {
		d.originalDrawOver(gc, rect)
	}
	if d.inDragOver {
		r := d.Table.ContentRect(false).Inset(unison.NewUniformInsets(1))
		paint := unison.ThemeWarning.Paint(gc, r, paintstyle.Stroke)
		paint.SetStrokeWidth(2)
		paint.SetColorFilter(unison.Alpha30Filter())
		gc.DrawRect(r, paint)
		paint.SetColorFilter(nil)
		paint.SetPathEffect(unison.DashEffect())
		gc.DrawLine(d.left, d.top, r.Right(), d.top, paint)
	}
}

func (d *TableDrop[T, U]) DataDragOverCallback(where unison.Point, data map[string]any) bool {
	if d.Table.filteredRows != nil {
		return false
	}
	var zero *Node[T]
	d.inDragOver = false
	if dd, ok := data[d.DragKey]; ok {
		if d.TableDragData, ok = dd.(*TableDragData[T]); ok {
			d.inDragOver = true
			last := d.Table.LastRowIndex()
			contentRect := d.Table.ContentRect(false)
			hierarchyColumnIndex := d.Table.ColumnIndexForID(d.Table.HierarchyColumnID)
			if where.Y >= contentRect.Bottom()-2 {

				d.TargetParent = zero
				d.TargetIndex = d.Table.RootRowCount()
				rect := d.Table.RowFrame(last)
				d.top = min(rect.Bottom()+1+d.Table.Padding.Bottom, contentRect.Bottom()-1)
				d.left, _ = d.Table.ColumnEdges(max(hierarchyColumnIndex, 0))
				d.Table.MarkForRedraw()
				return true
			}
			if rowIndex := d.Table.OverRow(where.Y); rowIndex != -1 {

				d.TargetIndex = -1
				row := d.Table.RowFromIndex(rowIndex)
				rect := d.Table.CellFrame(rowIndex, max(hierarchyColumnIndex, 0))
				if where.Y >= d.Table.RowFrame(rowIndex).CenterY() {
					d.top = min(rect.Bottom()+1+d.Table.Padding.Bottom, contentRect.Bottom()-1)
					d.left = rect.X

					if row.CanHaveChildren() {

						d.TargetParent = row
						d.TargetIndex = 0
						if hierarchyColumnIndex != -1 {
							d.left += d.Table.HierarchyIndent
						}
					} else {

						d.TargetParent = row.Parent()
						if row = d.Table.RowFromIndex(rowIndex + 1); row == zero {
							if d.TargetParent == zero {
								d.TargetIndex = len(d.Table.RootRows())
							} else {
								d.TargetIndex = len(d.TargetParent.Children)
							}
						}
					}
				} else {

					d.TargetParent = row.Parent()
					d.top = max(rect.Y-d.Table.Padding.Bottom, 1)
					d.left = rect.X
				}
				if d.TargetIndex == -1 && row != zero {
					var children []*Node[T]
					if d.TargetParent == zero {
						children = d.Table.RootRows()
					} else {
						children = d.TargetParent.Children
					}
					for i, child := range children {
						if child.UUID() == row.UUID() {
							d.TargetIndex = i
							break
						}
					}
					if d.TargetIndex == -1 {
						d.TargetIndex = len(children)
					}
				}

				if d.TargetParent != zero && d.TableDragData.Table == d.Table {
					for _, r := range d.TableDragData.Rows {
						if RowContainsRow(r, d.TargetParent) {

							d.inDragOver = false
							d.TargetParent = zero
							break
						}
					}
				}
				d.Table.MarkForRedraw()
				return true
			}

			d.TargetParent = zero
			d.TargetIndex = d.Table.RootRowCount()
			rect := d.Table.RowFrame(last)
			d.top = min(rect.Bottom()+1+d.Table.Padding.Bottom, contentRect.Bottom()-1)
			d.left, _ = d.Table.ColumnEdges(max(hierarchyColumnIndex, 0))
			d.Table.MarkForRedraw()
			return true
		}
	}
	return false
}

func (d *TableDrop[T, U]) DataDragExitCallback() {
	d.inDragOver = false
	var zero *Node[T]
	d.TargetParent = zero
	d.Table.MarkForRedraw()
}

func (d *TableDrop[T, U]) DataDragDropCallback(_ unison.Point, data map[string]any) {
	var savedScrollX, savedScrollY float32
	if scroller := d.Table.ScrollRoot(); scroller != nil {
		savedScrollX, savedScrollY = scroller.Position()
		defer func() {
			scroller.SetPosition(savedScrollX, savedScrollY)
		}()
	}
	var zero *Node[T]
	d.inDragOver = false
	var ok bool
	if d.TableDragData, ok = data[d.DragKey].(*TableDragData[T]); ok {
		d.AllDragData = data

		move := d.shouldMoveDataCallback(d.TableDragData.Table, d.Table)
		var undo *unison.UndoEdit[U]
		if d.willDropCallback != nil {
			undo = d.willDropCallback(d.TableDragData.Table, d.Table, move)
		}
		rows := slices.Clone(d.TableDragData.Rows)
		if move {

			commonParents := collectCommonParents(rows)
			for parent, list := range commonParents {
				var children []*Node[T]
				if parent == zero {
					children = d.TableDragData.Table.RootRows()
				} else {
					children = parent.Children
				}
				list = d.pruneRows(parent, children, makeRowSet(list))
				if parent == zero {
					d.TableDragData.Table.SetRootRows(list)
				} else {
					parent.SetChildren(list)
				}
			}
			d.TableDragData.Table.ClearSelection()
			d.TableDragData.Table.SyncToModel()

			for _, row := range rows {
				row.SetParent(d.TargetParent)
			}

			if d.Table != d.TableDragData.Table {
				if d.Table != d.TableDragData.Table && d.TableDragData.Table.DragRemovedRowsCallback != nil {
					d.TableDragData.Table.DragRemovedRowsCallback()
				}
			}
		} else {
			for i, row := range rows {
				rows[i] = row.CloneForTarget(d.Table, d.TargetParent)
			}
		}

		var targetRows []*Node[T]
		if d.TargetParent == zero {
			targetRows = d.Table.RootRows()
		} else {
			targetRows = d.TargetParent.Children
		}
		targetRows = slices.Insert(slices.Clone(targetRows), max(min(d.TargetIndex, len(targetRows)), 0), rows...)
		if d.TargetParent == zero {
			d.Table.SetRootRows(targetRows)
		} else {
			d.TargetParent.SetChildren(targetRows)
			d.Table.SyncToModel()
		}

		selMap := make(map[uuid.UUID]bool, len(rows))
		for _, row := range rows {
			selMap[row.UUID()] = true
		}
		d.Table.SetSelectionMap(selMap)

		if d.Table.DropOccurredCallback != nil {
			d.Table.DropOccurredCallback()
		}

		if d.didDropCallback != nil {
			d.didDropCallback(undo, d.TableDragData.Table, d.Table, move)
		}
	}
	d.Table.MarkForRedraw()
	d.TargetParent = zero
	d.AllDragData = nil
	d.TableDragData = nil
}

func (d *TableDrop[T, U]) pruneRows(parent *Node[T], rows []*Node[T], movingSet map[uuid.UUID]bool) []*Node[T] {
	movingToThisParent := d.TargetParent == parent
	list := make([]*Node[T], 0, len(rows))
	for i, row := range rows {
		if movingSet[row.UUID()] {
			if movingToThisParent && d.TargetIndex >= i {
				d.TargetIndex--
			}
		} else {
			list = append(list, row)
		}
	}
	return list
}

func makeRowSet[T any](rows []*Node[T]) map[uuid.UUID]bool {
	set := make(map[uuid.UUID]bool, len(rows))
	for _, row := range rows {
		set[row.UUID()] = true
	}
	return set
}

func collectCommonParents[T any](rows []*Node[T]) map[*Node[T]][]*Node[T] {
	m := make(map[*Node[T]][]*Node[T])
	for _, row := range rows {
		parent := row.Parent()
		m[parent] = append(m[parent], row)
	}
	return m
}
