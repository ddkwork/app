package widget

import (
	"time"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/unison"
	"github.com/ddkwork/unison/enums/align"
	"github.com/ddkwork/unison/enums/side"
)

func NewDockContainerClose(dock *TabClose) *unison.DockContainer {
	return unison.Ancestor[*unison.DockContainer](dock)
}

func NewDockContainer(dock *Tab) *unison.DockContainer {
	return unison.Ancestor[*unison.DockContainer](dock)
}

type (
	HSplitCloser struct {
		*unison.Dock
		LeftContainer  *unison.DockContainer
		RightContainer *unison.DockContainer
		pos            float32
		left, right    *TabClose // when add item,need setCurrent Dock
	}
	VSplitCloser struct {
		*unison.Dock
		TopContainer    *unison.DockContainer
		BottomContainer *unison.DockContainer
		pos             float32
		top, bottom     *TabClose
	}
	HSplit struct {
		*unison.Dock
		LeftContainer  *unison.DockContainer
		RightContainer *unison.DockContainer
		pos            float32
		left, right    *Tab
	}
	VSplit struct {
		*unison.Dock
		TopContainer    *unison.DockContainer
		BottomContainer *unison.DockContainer
		pos             float32
		top, bottom     *Tab
	}
)

func NewHSplitCloser(left, right *TabClose, scale float32) *HSplitCloser {
	mylog.Check(scale < 0.9) // 中心线往右或者往下的比例，0.8以上不可见就没必要拆分了，
	// 水平拆分：左半部分往右靠的比例，负小数则左半部分往左靠
	// 垂直拆分：上半部分往下靠的比例，负小数则上半部分往上靠
	dock := unison.NewDock()
	dock.AsPanel().SetLayoutData(&unison.FlexLayoutData{ // 没有它，tab设置的滚动条不显示
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	// SetScrollLayout(dock, 1)//todo test
	display := unison.PrimaryDisplay()
	r := display.Usable.Size
	pos := float32(0)

	dock.DockTo(left, nil, side.Left)
	LeftContainer := NewDockContainerClose(left)
	dock.DockTo(right, LeftContainer, side.Right)
	RightContainer := NewDockContainerClose(right)
	pos = (r.Width/(display.ScaleX*1000)/2 + scale) * 1000

	dock.RootDockLayout().SetDividerPosition(pos)
	unison.InvokeTaskAfter(func() { dock.RootDockLayout().SetDividerPosition(pos) }, time.Millisecond)
	s := &HSplitCloser{
		Dock:           dock,
		LeftContainer:  LeftContainer,
		RightContainer: RightContainer,
		pos:            pos,
		left:           left,
		right:          right,
	}
	s.SetCurrentDockable()
	return s
}

func (s *HSplitCloser) SetCurrentDockable() {
	s.LeftContainer.SetCurrentDockable(s.left)
	s.RightContainer.SetCurrentDockable(s.right)
}

func (s *HSplitCloser) AddLeftItem(tabClose *TabClose) {
	s.AddLeftItemAt(tabClose, -1)
}

func (s *HSplitCloser) AddLeftItemAt(tabClose *TabClose, index int) {
	s.LeftContainer.Stack(tabClose, index)
	s.SetCurrentDockable()
}

func (s *HSplitCloser) AddRightItem(tabClose *TabClose) {
	s.AddRightItemAt(tabClose, -1)
}

func (s *HSplitCloser) AddRightItemAt(tabClose *TabClose, index int) {
	s.RightContainer.Stack(tabClose, index)
	s.SetCurrentDockable()
}

func NewVSplitCloser(Top, Bottom *TabClose, scale float32) *VSplitCloser {
	mylog.Check(scale < 0.9) // 中心线往右或者往下的比例，0.8以上不可见就没必要拆分了，
	// 水平拆分：左半部分往右靠的比例，负小数则左半部分往左靠
	// 垂直拆分：上半部分往下靠的比例，负小数则上半部分往上靠

	dock := unison.NewDock()
	dock.AsPanel().SetLayoutData(&unison.FlexLayoutData{ // 没有它，tab设置的滚动条不显示
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	// SetScrollLayout(dock, 1)//todo test
	display := unison.PrimaryDisplay()
	r := display.Usable.Size
	pos := float32(0)
	dock.DockTo(Top, nil, side.Top)
	TopContainer := NewDockContainerClose(Top)
	dock.DockTo(Bottom, TopContainer, side.Bottom)
	BottomContainer := NewDockContainerClose(Bottom)
	pos = (1 - (r.Height/(display.ScaleX*1000)/2 + scale)) * 1000
	dock.RootDockLayout().SetDividerPosition(pos)
	unison.InvokeTaskAfter(func() { dock.RootDockLayout().SetDividerPosition(pos) }, time.Millisecond)
	s := &VSplitCloser{
		Dock:            dock,
		TopContainer:    TopContainer,
		BottomContainer: BottomContainer,
		pos:             pos,
		top:             Top,
		bottom:          Bottom,
	}
	s.SetCurrentDockable()
	return s
}

func (s *VSplitCloser) SetCurrentDockable() {
	s.TopContainer.SetCurrentDockable(s.top)
	s.BottomContainer.SetCurrentDockable(s.bottom)
}

func (s *VSplitCloser) AddTopItem(tabClose *TabClose) {
	s.AddTopItemAt(tabClose, -1)
}

func (s *VSplitCloser) AddTopItemAt(tabClose *TabClose, index int) {
	s.TopContainer.Stack(tabClose, index)
	s.SetCurrentDockable()
}

func (s *VSplitCloser) AddBottomItem(tabClose *TabClose) {
	s.AddBottomItemAt(tabClose, -1)
}

func (s *VSplitCloser) AddBottomItemAt(tabClose *TabClose, index int) {
	s.BottomContainer.Stack(tabClose, index)
	s.SetCurrentDockable()
}

// ////////////////////////////////////////////
func NewHSplit(left, right *Tab, scale float32) *HSplit {
	mylog.Check(scale < 0.9) // 中心线往右或者往下的比例，0.8以上不可见就没必要拆分了，
	// 水平拆分：左半部分往右靠的比例，负小数则左半部分往左靠
	// 垂直拆分：上半部分往下靠的比例，负小数则上半部分往上靠

	dock := unison.NewDock()
	dock.AsPanel().SetLayoutData(&unison.FlexLayoutData{ // 没有它，tab设置的滚动条不显示
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	// SetScrollLayout(dock, 1) //todo test
	display := unison.PrimaryDisplay()
	r := display.Usable.Size
	pos := float32(0)

	dock.DockTo(left, nil, side.Left)
	LeftContainer := NewDockContainer(left)
	dock.DockTo(right, LeftContainer, side.Right)
	RightContainer := NewDockContainer(right)
	pos = (r.Width/(display.ScaleX*1000)/2 + scale) * 1000

	dock.RootDockLayout().SetDividerPosition(pos)
	unison.InvokeTaskAfter(func() { dock.RootDockLayout().SetDividerPosition(pos) }, time.Millisecond)
	s := &HSplit{
		Dock:           dock,
		LeftContainer:  LeftContainer,
		RightContainer: RightContainer,
		pos:            pos,
		left:           left,
		right:          right,
	}
	s.SetCurrentDockable()
	return s
}

func (s *HSplit) SetCurrentDockable() {
	s.LeftContainer.SetCurrentDockable(s.left)
	s.RightContainer.SetCurrentDockable(s.right)
}

func (s *HSplit) AddLeftItem(Tab *Tab) {
	s.AddLeftItemAt(Tab, -1)
}

func (s *HSplit) AddLeftItemAt(Tab *Tab, index int) {
	s.LeftContainer.Stack(Tab, index)
	s.SetCurrentDockable()
}

func (s *HSplit) AddRightItem(Tab *Tab) {
	s.AddRightItemAt(Tab, -1)
}

func (s *HSplit) AddRightItemAt(Tab *Tab, index int) {
	s.RightContainer.Stack(Tab, index)
	s.SetCurrentDockable()
}

func NewVSplit(Top, Bottom *Tab, scale float32) *VSplit {
	mylog.Check(scale < 0.9) // 中心线往右或者往下的比例，0.8以上不可见就没必要拆分了，
	// 水平拆分：左半部分往右靠的比例，负小数则左半部分往左靠
	// 垂直拆分：上半部分往下靠的比例，负小数则上半部分往上靠

	dock := unison.NewDock()
	dock.AsPanel().SetLayoutData(&unison.FlexLayoutData{ // 没有它，tab设置的滚动条不显示
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	// SetScrollLayout(dock, 1)//todo test
	display := unison.PrimaryDisplay()
	r := display.Usable.Size
	pos := float32(0)
	dock.DockTo(Top, nil, side.Top)
	TopContainer := NewDockContainer(Top)
	dock.DockTo(Bottom, TopContainer, side.Bottom)
	BottomContainer := NewDockContainer(Bottom)
	pos = (1 - (r.Height/(display.ScaleX*1000)/2 + scale)) * 1000
	dock.RootDockLayout().SetDividerPosition(pos)
	unison.InvokeTaskAfter(func() { dock.RootDockLayout().SetDividerPosition(pos) }, time.Millisecond)
	s := &VSplit{
		Dock:            dock,
		TopContainer:    TopContainer,
		BottomContainer: BottomContainer,
		pos:             pos,
		top:             Top,
		bottom:          Bottom,
	}
	s.SetCurrentDockable()
	return s
}

func (s *VSplit) SetCurrentDockable() {
	s.TopContainer.SetCurrentDockable(s.top)
	s.BottomContainer.SetCurrentDockable(s.bottom)
}

func (s *VSplit) AddTopItem(Tab *Tab) {
	s.AddTopItemAt(Tab, -1)
}

func (s *VSplit) AddTopItemAt(Tab *Tab, index int) {
	s.TopContainer.Stack(Tab, index)
	s.SetCurrentDockable()
}

func (s *VSplit) AddBottomItem(Tab *Tab) {
	s.AddBottomItemAt(Tab, -1)
}

func (s *VSplit) AddBottomItemAt(Tab *Tab, index int) {
	s.BottomContainer.Stack(Tab, index)
	s.SetCurrentDockable()
}
