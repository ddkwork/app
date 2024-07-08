package widget

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"unicode"

	"github.com/richardwilkes/toolbox/i18n"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
)

type API interface {
	Layout() unison.Paneler
}

// todo
// D:\workspace\workspace\app\widget\chapar
// D:\workspace\workspace\app\widget\corewidget
// D:\workspace\workspace\app\widget\widget
// 三个包的tui树形表格移除之前测试本包的树形表格单元测试，同时需要再gio移植完成之后
// 高优先级：
// 代码编辑器，为mitmproxy写一次
// list控件实现一个todo，支持拖动随时调整优先级
// 垂直拆分和水平拆分控件
// 结构体树形表格，完成反射node填充
// explorer
// 结构体显示控件布局
// asmView的着色和层级控制

//func newWidget[T any]() Widget[T] {
//	return &widget[T]{}
//}

func LabelStyle(label *unison.Label) {
	if label.String() != "" || label.Drawable != nil {
		label.MouseEnterCallback = func(_ unison.Point, _ unison.Modifiers) bool {
			size := label.Font.Size()
			size++
			label.Font = label.Font.Face().Font(size)
			label.MarkForRedraw()
			return true
		}
		label.MouseExitCallback = func() bool {
			size := label.Font.Size()
			size--
			label.Font = label.Font.Face().Font(size)
			label.MarkForRedraw()
			return true
		}
		label.MouseDownCallback = func(where unison.Point, _, _ int, _ unison.Modifiers) bool {
			return false
		}
	}
}

type Panel struct {
	unison.Panel
}

func NewPanel() *Panel {
	p := &Panel{
		Panel: unison.Panel{},
	}
	p.Self = p
	SetScrollLayout(p, 1)
	return p
}

func PanelSetBorder(panel unison.Paneler) {
	panel.AsPanel().SetBorder(unison.NewEmptyBorder(unison.NewSymmetricInsets(unison.StdHSpacing, unison.StdVSpacing)))
}

func (p *Panel) AddChildren(children ...unison.Paneler) *Panel {
	for _, child := range children {
		p.AddChild(child)
	}
	return p
}

func (p *Panel) SetFlexLayout(layout unison.FlexLayout) *Panel {
	p.SetLayout(&layout)
	return p
}

func SetScrollLayout(paneler unison.Paneler, Columns int) {
	paneler.AsPanel().SetLayout(&unison.FlexLayout{Columns: Columns})
	paneler.AsPanel().SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		// HGrab:  true,
		// VGrab:  true,
	})
}

func NewScrollPanelFill(content unison.Paneler) *unison.ScrollPanel {
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(content, behavior.Fill, behavior.Fill) // 滚动条与布局边缘重叠
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.Content().AsPanel().ValidateScrollRoot()
	return scrollArea
}

func NewScrollPanelHintedFill(content unison.Paneler) *unison.ScrollPanel {
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(content, behavior.HintedFill, behavior.Fill) // 滚动条在布局之外
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.Content().AsPanel().ValidateScrollRoot()
	return scrollArea
}

func NewToolBar(buttons ...*unison.Button) unison.Paneler {
	panel := unison.NewPanel()
	PanelSetBorder(panel)
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	field := unison.NewField()
	field.Tooltip = unison.NewTooltipWithText("global filter...")
	field.MinimumTextWidth = 100
	field.SetLayoutData(align.Middle)
	field.EditableInk = unison.RGB(43, 43, 43)
	panel.AddChild(field)
	for _, button := range buttons {
		button.CornerRadius = 4
		button.EdgeInk = unison.ThemeSurfaceEdge
		button.SetLayoutData(align.Middle)
		panel.AddChild(button)
	}
	return panel
}

func createImageButton(img *unison.Image, actionText string, panel *unison.Panel) *unison.Button {
	btn := unison.NewButton()
	btn.Drawable = img
	btn.ClickCallback = func() { slog.Info(actionText) }
	btn.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("Tooltip for: %s", actionText))
	btn.SetLayoutData(align.Middle)
	panel.AddChild(btn)
	return btn
}

func NewSeparator() *unison.Separator {
	sep := unison.NewSeparator()
	sep.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
	})
	return sep
}

func NewImageButton[T stream.Type](tooltip string, imageBuf T, clickCallback func()) *unison.Button {
	button := NewButton("", clickCallback)
	button.Tooltip = unison.NewTooltipWithText(tooltip)
	fromBytes := mylog.Check2(unison.NewImageFromBytes(stream.NewBuffer(imageBuf).Bytes(), 0.5))
	button.Drawable = &unison.SizedDrawable{
		Drawable: fromBytes,
		Size:     unison.NewSize(24, 24),
	}
	return button
}

func NewButton(Text string, ClickCallback func()) *unison.Button {
	b := unison.NewButton()
	b.ClickCallback = ClickCallback
	b.CornerRadius = 18
	b.HMargin = 12
	b.SetTitle(Text)
	b.EdgeInk = unison.White

	b.MouseEnterCallback = func(where unison.Point, mod unison.Modifiers) bool { // todo bug new version not working
		size := b.Font.Size()
		size++
		b.Font = b.Font.Face().Font(size)
		b.SetTitle(Text)
		b.MarkForRedraw()
		return true
	}

	b.MouseExitCallback = func() bool {
		size := b.Font.Size()
		size--
		b.Font = b.Font.Face().Font(size)
		b.SetTitle(Text)
		b.MarkForRedraw()
		return true
	}

	return b
}

func CreatePopupMenu(parent *unison.Panel, p *unison.PopupMenu[string], selectIndex int, tooltip string, titles ...string) *unison.PopupMenu[string] {
	p.Tooltip = unison.NewTooltipWithText(tooltip)
	for _, title := range titles {
		if title == "" {
			p.AddSeparator()
		} else {
			p.AddItem(title)
		}
	}
	p.SelectIndex(selectIndex)
	parent.AddChild(p)
	return p
}

func (s *StructView[T]) getFieldValues() []string {
	values := make([]string, len(s.Editors))
	for i, field := range s.Editors {
		values[i] = field.Field.Text()
	}
	return values
}

func (s *StructView[T]) Unmarshal(fnValues func(values []string)) {
	fnValues(s.getFieldValues())
}

func (s *StructView[T]) UpdateField(index int, value string) {
	s.Editors[index].Field.SetText(value)
}
func (s *StructView[T]) Update(data T) {
	s.MetaData = data
}

func NewLogView() *unison.Field {
	f := unison.NewMultiLineField()
	f.MinimumTextWidth = 666
	f.SetText(`
log ...












`)
	return f
}

type (
	StructView[T any] struct {
		unison.Panel
		MetaData T
		Editors  []StructEditor
	}
	StructEditor struct {
		*unison.Label
		*unison.Field
		KeyValueToolTip
	}
	structField struct { // for  reflect.VisibleFields
		reflect.StructField
		KeyValueToolTip
	}
	KeyValueToolTip struct {
		Key     string
		Value   string
		Tooltip string
	}

	RowValueType interface {
		*unison.PopupMenu[string] | []*unison.Button | *unison.Field //|constraints.Ordered
	}
)

type StructViewPanel struct {
	unison.Panel
	keyValuePanel *unison.Panel
}

func NewStructViewPanel() *StructViewPanel {
	s := &StructViewPanel{}
	s.Self = s
	s.SetLayout(&unison.FlexLayout{Columns: 1})
	s.AddChild(NewVSpacer())
	s.keyValuePanel = NewKeyValuePanel()
	s.AddChild(s.keyValuePanel)
	return s
}

func AddRowForStructViewPanel[T RowValueType](s *StructViewPanel, key, tooltip string, value T) {
	s.keyValuePanel.AddChild(NewLabelRightAlign(KeyValueToolTip{
		Key:     key,
		Value:   "",
		Tooltip: tooltip,
	}))
	switch v := any(value).(type) {
	case *unison.PopupMenu[string]:
		v.SelectIndex(0)
		s.keyValuePanel.AddChild(v)
	case []*unison.Button:
		buttonCount := len(v)
		buttonPanel := NewPanel().SetFlexLayout(unison.FlexLayout{
			Columns:      buttonCount + 1,
			HSpacing:     unison.StdHSpacing * 2,
			VSpacing:     unison.StdVSpacing,
			EqualColumns: true,
		})
		buttonPanel.AddChild(unison.NewPanel()) // left spacer
		for _, button := range v {
			buttonPanel.AddChild(button)
		}
		buttonPanel.SetLayoutData(&unison.FlexLayoutData{
			HSpan:  3,
			VSpan:  1,
			HAlign: align.End,
			VAlign: align.Middle,
		})
		s.keyValuePanel.AddChild(buttonPanel)
	case *unison.Field:
		v.MinimumTextWidth = 300
		s.keyValuePanel.AddChild(v)
	}
}

func NewStructView[T any](data T, marshal func(data T) (values []CellData)) (view *StructView[T], kvPanel *unison.Panel) {
	visibleFields := stream.ReflectVisibleFields(data)
	fields := make([]structField, len(visibleFields))
	values := marshal(data)
	if len(visibleFields) != len(values) {
		mylog.Check("NewStructView init error : len(visibleFields) != len(values)")
	}
	for i, field := range visibleFields {
		f := structField{
			StructField: field,
			KeyValueToolTip: KeyValueToolTip{
				Key:     field.Name,
				Value:   values[i].Text,
				Tooltip: "",
			},
		}
		f.Tooltip = f.SetTooltip()
		fields[i] = f
	}
	view = &StructView[T]{
		Panel:    unison.Panel{},
		MetaData: data,
		Editors:  make([]StructEditor, 0, len(visibleFields)),
	}
	view.Self = view
	// view.undoMgr = unison.NewUndoManager(100, func(err error) { errs.Log(err) })
	view.SetLayout(&unison.FlexLayout{Columns: 1})
	view.AddChild(NewVSpacer())
	kvPanel = NewKeyValuePanel()
	for _, editor := range fields {
		key := NewLabelRightAlign(editor.KeyValueToolTip)
		value := NewFieldLeftAlign(editor.KeyValueToolTip)
		view.Editors = append(view.Editors, StructEditor{
			Label:           key,
			Field:           value,
			KeyValueToolTip: editor.KeyValueToolTip,
		})
		kvPanel.AddChild(key)
		kvPanel.AddChild(value)
	}
	view.AddChild(kvPanel)
	kvPanel.AddChild(NewVSpacer())
	// NewScrollPanelHintedFill(view.AsPanel())
	return
}

func (s structField) SetTooltip() string {
	b := stream.NewBuffer("")
	fnSep := func() { b.WriteString("    |    ") }

	b.WriteString(s.Name)
	fnSep()

	b.WriteString(s.Type.String())
	fnSep()

	get := s.Tag.Get("json")
	if get != "" {
		b.WriteString(get)
		fnSep()
	}

	b.WriteString(s.Value)
	return b.String()
}

func NewFieldContextMenuItems(filed *unison.Field) {
	NewContextMenuItems(filed, filed.DefaultMouseDown,
		ContextMenuItem{
			Title: "Copy",
			Can:   func(any) bool { return filed.CanCopy() },
			Do:    func(a any) { filed.Copy() },
		},
		ContextMenuItem{
			Title: "Paste",
			Can:   func(any) bool { return filed.CanPaste() },
			Do:    func(a any) { filed.Paste() },
		},
		ContextMenuItem{
			Title: "Cut",
			Can:   func(any) bool { return filed.CanCut() },
			Do:    func(a any) { filed.Cut() },
		},
		ContextMenuItem{
			Title: "Delete",
			Can:   func(any) bool { return filed.CanDelete() },
			Do:    func(a any) { filed.Delete() },
		},
		ContextMenuItem{
			Title: "SelectAll",
			Can:   func(any) bool { return filed.CanSelectAll() },
			Do:    func(a any) { filed.SelectAll() },
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
}

func NewApplyCancelButtonPanel(parent *unison.Panel, applyCallback, cancelCallback func()) {
	NewVSpacer()
	buttonPanel := NewPanel().SetFlexLayout(unison.FlexLayout{
		Columns:      3,
		HSpacing:     unison.StdHSpacing * 2,
		VSpacing:     unison.StdVSpacing,
		EqualColumns: true,
	})
	buttonPanel.AddChild(unison.NewPanel()) // left spacer

	applyButton := NewButton("apply", applyCallback)
	cancelButton := NewButton("cancel", cancelCallback)
	buttonPanel.AddChild(applyButton)
	buttonPanel.AddChild(cancelButton)
	buttonPanel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  3,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	parent.AddChild(buttonPanel)
}

func NewButtonsPanel(titles []string, callbacks ...func()) *unison.Panel {
	buttonCount := len(titles)
	callbackCount := len(callbacks)
	mylog.Check(buttonCount == callbackCount)
	buttonPanel := NewPanel().SetFlexLayout(unison.FlexLayout{
		Columns:  buttonCount + 1,
		HSpacing: unison.StdHSpacing * 2,
		VSpacing: unison.StdVSpacing,
	})
	buttonPanel.AddChild(unison.NewPanel()) // left spacer

	for i, title := range titles {
		button := NewButton(title, nil)
		if i < len(callbacks) {
			button.ClickCallback = callbacks[i]
		}
		buttonPanel.AddChild(button)
	}
	buttonPanel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  3,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	buttonPanel.AddChild(NewVSpacer())
	return buttonPanel.AsPanel()
}

func NewVSpacer() *unison.Panel {
	vSpacer := unison.NewPanel()
	vSpacer.AsPanel().SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(10)))
	vSpacer.AsPanel().SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: 10,
	})
	return vSpacer
}

func NewFieldLeftAlign(kvt KeyValueToolTip) *unison.Field {
	field := unison.NewField()
	field.MinimumTextWidth = 520
	NewFieldContextMenuItems(field)
	field.SetText(kvt.Value)
	field.SetLayoutData(&unison.FlexLayoutData{
		SizeHint: unison.Size{},
		MinSize:  unison.Size{},
		HSpan:    0,
		VSpan:    0,
		HAlign:   align.Start,
		HGrab:    true,
		VGrab:    true,
	})
	return field
}

func NewKeyValuePanel() *unison.Panel {
	kvPanel := unison.NewPanel()
	kvPanel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		HAlign:   0,
		VAlign:   0,
	})
	return kvPanel
}

func NewLabelRightAlign(kvt KeyValueToolTip) *unison.Label {
	label := unison.NewLabel()
	label.SetTitle(i18n.Text(kvt.Key))
	label.Tooltip = unison.NewTooltipWithText(kvt.Tooltip)
	label.SetLayoutData(&unison.FlexLayoutData{
		SizeHint: unison.Size{},
		HSpan:    0,
		VSpan:    0,
		HAlign:   align.End,    // 右对齐，not used
		VAlign:   align.Middle, // 垂直居中
		HGrab:    true,
		VGrab:    true,
	})
	label.HAlign = align.End
	label.VAlign = align.Middle
	return label
}

// //////////////////////
func createFieldsAndListPanel() *unison.Panel {
	// Create a wrapper to put them side-by-side
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:      2,
		HSpacing:     10,
		VSpacing:     unison.StdVSpacing,
		EqualColumns: true,
	})

	// Add the text fields to the left side
	textFieldsPanel := createTextFieldsPanel()
	textFieldsPanel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	wrapper.AddChild(textFieldsPanel)

	// Add the list to the right side
	wrapper.AddChild(createListPanel())

	return wrapper
}

func createTextFieldsPanel() *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	createTextField("Field 1:", "First Text Field", panel)
	createTextField("Field 2:", "Second Text Field (disabled)", panel).SetEnabled(false)
	field := createTextField("Longer Label:", "", panel)
	field.Watermark = "Password Field"
	field.ObscurementRune = '●'
	field = createTextField("Field 4:", "", panel)
	field.HAlign = align.End
	field.Watermark = "Enter only numbers"
	field.ValidateCallback = func() bool {
		for _, r := range field.Text() {
			if !unicode.IsDigit(r) {
				return false
			}
		}
		return true
	}
	createMultiLineTextField("Field 5:", "One\nTwo\nThree", panel)
	return panel
}

func createTextField(labelText, fieldText string, panel *unison.Panel) *unison.Field {
	lbl := unison.NewLabel()
	lbl.SetTitle(labelText)
	lbl.HAlign = align.End
	lbl.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	panel.AddChild(lbl)
	field := unison.NewField()
	field.SetText(fieldText)
	field.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for %v", field))
	panel.AddChild(field)
	return field
}

func createMultiLineTextField(labelText, fieldText string, panel *unison.Panel) *unison.Field {
	lbl := unison.NewLabel()
	lbl.SetTitle(labelText)
	lbl.HAlign = align.End
	lbl.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	panel.AddChild(lbl)
	field := unison.NewMultiLineField()
	field.SetText(fieldText)
	field.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for %v", field))
	panel.AddChild(field)
	return field
}

func createListPanel() *unison.Panel {
	lst := unison.NewList[string]()
	lst.Append(
		"One",
		"Two",
		"Three with some long text to make it interesting",
		"Four",
		"Five",
	)
	lst.NewSelectionCallback = func() {
		var buffer strings.Builder
		buffer.WriteString("Selection changed in the list. Now:")
		index := -1
		first := true
		for {
			index = lst.Selection.NextSet(index + 1)
			if index == -1 {
				break
			}
			if first {
				first = false
			} else {
				buffer.WriteString(",")
			}
			fmt.Fprintf(&buffer, " %d", index)
		}
		slog.Info(buffer.String())
	}
	lst.DoubleClickCallback = func() {
		slog.Info("Double-clicked on the list")
	}
	_, prefSize, _ := lst.Sizes(unison.Size{})
	lst.SetFrameRect(unison.Rect{Size: prefSize})
	scroller := unison.NewScrollPanel()
	scroller.SetBorder(unison.NewLineBorder(unison.ThemeSurfaceEdge, 0, unison.NewUniformInsets(1), false))
	scroller.SetContent(lst, behavior.Fill, behavior.Fill)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	return scroller.AsPanel()
}
