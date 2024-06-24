package widget

import (
	"fmt"

	"github.com/richardwilkes/unison"

	"github.com/richardwilkes/unison/enums/paintstyle"
)

type dragDrawable struct {
	label *unison.Label
}

func NewTableDragDrawable[T any](data *TableDragData[T], svg *unison.SVG, singularName, pluralName string) unison.Drawable {
	label := unison.NewLabel()
	label.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		r := rect.Inset(unison.NewUniformInsets(1))
		corner := r.Height / 2
		gc.SaveWithOpacity(0.7)
		gc.DrawRoundedRect(r, corner, corner, data.Table.SelectionInk.Paint(gc, r, paintstyle.Fill))
		gc.DrawRoundedRect(r, corner, corner, data.Table.OnSelectionInk.Paint(gc, r, paintstyle.Stroke))
		gc.Restore()
		label.DefaultDraw(gc, rect)
	}
	label.OnBackgroundInk = data.Table.OnSelectionInk
	label.SetBorder(unison.NewEmptyBorder(unison.Insets{
		Top:    4,
		Left:   label.Font.LineHeight(),
		Bottom: 4,
		Right:  label.Font.LineHeight(),
	}))
	if count := CountTableRows(data.Rows); count == 1 {
		label.SetTitle(fmt.Sprintf("1 %s", singularName))
	} else {
		label.SetTitle(fmt.Sprintf("%d %s", count, pluralName))
	}
	if svg != nil {
		baseline := label.Font.Baseline()
		label.Drawable = &unison.DrawableSVG{
			SVG:  svg,
			Size: unison.Size{Width: baseline, Height: baseline},
		}
	}
	_, pref, _ := label.Sizes(unison.Size{})
	label.SetFrameRect(unison.Rect{Size: pref})
	return &dragDrawable{label: label}
}

func (d *dragDrawable) LogicalSize() unison.Size {
	return d.label.FrameRect().Size
}

func (d *dragDrawable) DrawInRect(canvas *unison.Canvas, rect unison.Rect, _ *unison.SamplingOptions, _ *unison.Paint) {
	d.label.Draw(canvas, rect)
}
