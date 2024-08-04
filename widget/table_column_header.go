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
	"github.com/ddkwork/unison/enums/align"
	"github.com/ddkwork/unison/enums/paintstyle"
	"github.com/ddkwork/unison/enums/side"
)

// TableColumnHeader defines the methods a table column header must implement.
type TableColumnHeader[T any] interface {
	unison.Paneler
	SortState() unison.SortState
	SetSortState(state unison.SortState)
}

// DefaultTableColumnHeaderTheme holds the default TableColumnHeaderTheme values for TableColumnHeaders. Modifying this
// data will not alter existing TableColumnHeaders, but will alter any TableColumnHeaders created in the future.
var DefaultTableColumnHeaderTheme = unison.LabelTheme{
	Font:            unison.LabelFont,
	OnBackgroundInk: unison.OnBackgroundColor,
	Gap:             3,
	HAlign:          align.Middle,
	VAlign:          align.Middle,
	Side:            side.Left,
}

// DefaultTableColumnHeader provides a default table column header panel.
type DefaultTableColumnHeader[T any] struct {
	unison.Label
	sortState     unison.SortState
	sortIndicator *unison.DrawableSVG
}

// NewTableColumnHeader creates a new table column header panel.
func NewTableColumnHeader[T any](title, tooltip string) *DefaultTableColumnHeader[T] {
	h := &DefaultTableColumnHeader[T]{
		Label: unison.Label{
			LabelTheme: DefaultTableColumnHeaderTheme,
			Text:       title,
		},
		sortState: unison.SortState{
			Order:     -1,
			Ascending: true,
			Sortable:  true,
		},
	}
	h.Self = h
	h.SetSizer(h.DefaultSizes)
	h.DrawCallback = h.DefaultDraw
	h.MouseUpCallback = h.DefaultMouseUp
	if tooltip != "" {
		h.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return h
}

// DefaultSizes provides the default sizing.
func (h *DefaultTableColumnHeader[T]) DefaultSizes(hint unison.Size) (minSize, prefSize, maxSize unison.Size) {
	prefSize = unison.LabelSize(h.TextCache.Text(h.Text, h.Font), h.Drawable, h.Side, h.Gap)

	// Account for the potential sort indicator
	baseline := h.Font.Baseline()
	prefSize.Width += h.LabelTheme.Gap + baseline
	if prefSize.Height < baseline {
		prefSize.Height = baseline
	}

	if b := h.Border(); b != nil {
		prefSize.AddInsets(b.Insets())
	}
	prefSize.GrowToInteger()
	prefSize.ConstrainForHint(hint)
	return prefSize, prefSize, prefSize
}

// DefaultDraw provides the default drawing.
func (h *DefaultTableColumnHeader[T]) DefaultDraw(canvas *unison.Canvas, _ unison.Rect) {
	r := h.ContentRect(false)
	if h.sortIndicator != nil {
		r.Width -= h.LabelTheme.Gap + h.sortIndicator.LogicalSize().Width
	}
	unison.DrawLabel(canvas, r, h.HAlign, h.VAlign, h.TextCache.Text(h.Text, h.Font), h.OnBackgroundInk, h.Drawable, h.Side,
		h.Gap, !h.Enabled())
	if h.sortIndicator != nil {
		size := h.sortIndicator.LogicalSize()
		r.X = r.Right() + h.LabelTheme.Gap
		r.Y += (r.Height - size.Height) / 2
		r.Size = size
		paint := h.OnBackgroundInk.Paint(canvas, r, paintstyle.Fill)
		if !h.Enabled() {
			paint.SetColorFilter(unison.Grayscale30Filter())
		}
		h.sortIndicator.DrawInRect(canvas, r, nil, paint)
	}
}

// SortState returns the current SortState.
func (h *DefaultTableColumnHeader[T]) SortState() unison.SortState {
	return h.sortState
}

// SetSortState sets the SortState.
func (h *DefaultTableColumnHeader[T]) SetSortState(state unison.SortState) {
	if h.sortState != state {
		h.sortState = state
		if h.sortState.Sortable && h.sortState.Order == 0 {
			baseline := h.Font.Baseline()
			if h.sortState.Ascending {
				h.sortIndicator = &unison.DrawableSVG{
					SVG:  unison.SortAscendingSVG,
					Size: unison.NewSize(baseline, baseline),
				}
			} else {
				h.sortIndicator = &unison.DrawableSVG{
					SVG:  unison.SortDescendingSVG,
					Size: unison.NewSize(baseline, baseline),
				}
			}
		} else {
			h.sortIndicator = nil
		}
		h.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (h *DefaultTableColumnHeader[T]) DefaultMouseUp(where unison.Point, _ int, _ unison.Modifiers) bool {
	if h.sortState.Sortable && h.ContentRect(false).ContainsPoint(where) {
		if header, ok := h.Parent().Self.(*TableHeader[T]); ok {
			header.SortOn(h)
			header.ApplySort()
		}
	}
	return true
}
