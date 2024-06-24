package widget

import (
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/side"
)

type TableColumnHeader[T any] interface {
	unison.Paneler
	SortState() unison.SortState
	SetSortState(state unison.SortState)
}

var DefaultTableColumnHeaderTheme = unison.LabelTheme{
	TextDecoration: unison.TextDecoration{
		Font:            unison.LabelFont,
		OnBackgroundInk: unison.ThemeOnSurface,
	},
	Gap:    3,
	HAlign: align.Middle,
	VAlign: align.Middle,
	Side:   side.Left,
}

type DefaultTableColumnHeader[T any] struct {
	*unison.Label
	sortState     unison.SortState
	sortIndicator *unison.DrawableSVG
}

func NewTableColumnHeader[T any](title, tooltip string) *DefaultTableColumnHeader[T] {
	h := &DefaultTableColumnHeader[T]{
		Label: unison.NewLabel(),
		sortState: unison.SortState{
			Order:     -1,
			Ascending: true,
			Sortable:  true,
		},
	}
	h.Self = h
	h.LabelTheme = DefaultTableColumnHeaderTheme
	h.SetTitle(title)
	h.SetSizer(h.DefaultSizes)
	h.DrawCallback = h.DefaultDraw
	h.MouseUpCallback = h.DefaultMouseUp
	if tooltip != "" {
		h.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return h
}

func (h *DefaultTableColumnHeader[T]) DefaultSizes(hint unison.Size) (minSize, prefSize, maxSize unison.Size) {
	prefSize, _ = unison.LabelContentSizes(h.Text, h.Drawable, h.Font, h.Side, h.Gap)

	baseline := h.Font.Baseline()
	prefSize.Width += h.LabelTheme.Gap + baseline
	if prefSize.Height < baseline {
		prefSize.Height = baseline
	}

	if b := h.Border(); b != nil {
		prefSize = prefSize.Add(b.Insets().Size())
	}
	prefSize = prefSize.Ceil().ConstrainForHint(hint)
	return prefSize, prefSize, prefSize
}

func (h *DefaultTableColumnHeader[T]) DefaultDraw(canvas *unison.Canvas, _ unison.Rect) {
	r := h.ContentRect(false)
	if h.sortIndicator != nil {
		r.Width -= h.LabelTheme.Gap + h.sortIndicator.LogicalSize().Width
	}
	unison.DrawLabel(canvas, r, h.HAlign, h.VAlign, h.Font, h.Text, h.OnBackgroundInk, nil, h.Drawable, h.Side, h.Gap,
		!h.Enabled())
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

func (h *DefaultTableColumnHeader[T]) SortState() unison.SortState {
	return h.sortState
}

func (h *DefaultTableColumnHeader[T]) SetSortState(state unison.SortState) {
	if h.sortState != state {
		h.sortState = state
		if h.sortState.Sortable && h.sortState.Order == 0 {
			baseline := h.Font.Baseline()
			if h.sortState.Ascending {
				h.sortIndicator = &unison.DrawableSVG{
					SVG:  unison.SortAscendingSVG,
					Size: unison.Size{Width: baseline, Height: baseline},
				}
			} else {
				h.sortIndicator = &unison.DrawableSVG{
					SVG:  unison.SortDescendingSVG,
					Size: unison.Size{Width: baseline, Height: baseline},
				}
			}
		} else {
			h.sortIndicator = nil
		}
		h.MarkForRedraw()
	}
}

func (h *DefaultTableColumnHeader[T]) DefaultMouseUp(where unison.Point, button int, _ unison.Modifiers) bool {
	if button == unison.ButtonRight {
		return false
	}
	if h.sortState.Sortable && where.In(h.ContentRect(false)) {
		if header, ok := h.Parent().Self.(*TableHeader[T]); ok {
			header.SortOn(h)
			header.ApplySort()
		}
	}
	return true
}
