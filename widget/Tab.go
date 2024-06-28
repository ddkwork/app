package widget

import (
	"github.com/richardwilkes/toolbox/i18n"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison"
)

var (
	_ unison.Dockable  = &Tab{}
	_ unison.Dockable  = &TabClose{}
	_ unison.TabCloser = &TabClose{}
)

func (t *TabClose) SetContent(content unison.Paneler) {
	content.AsPanel().RemoveFromParent()
	t.ScrollPanel = NewScrollPanelFill(t)
	t.MarkForLayoutAndRedraw()
}

type TabClose struct {
	unison.Panel
	title       string
	tooltip     string
	closeable   bool
	ScrollPanel *unison.ScrollPanel
}

func (t *TabClose) SetTitle(title string) {
	t.title = title
}

func (t *TabClose) SetTooltip(tooltip string) {
	t.tooltip = tooltip
}

func (t *TabClose) SetCloseable(closeable bool) {
	t.closeable = closeable
}

func (t *TabClose) SetScrollPanel(ScrollPanel *unison.ScrollPanel) {
	t.ScrollPanel = ScrollPanel
}

func NewTabCloseWithTable[T any](table *Node[T], header *TableHeader[T], title string, tooltip string, closeable bool) *TabClose {
	panel := NewPanel()
	panel.AddChild(table)
	panel.AddChild(header)

	tabClose := NewTabClose(title, tooltip, closeable, panel)
	tabClose.ScrollPanel.SetColumnHeader(header)
	return tabClose
}

type TabContent struct {
	Title     string
	Tooltip   string
	Closeable bool
	Panel     unison.Paneler
}

func NewTabCloses(tabContents ...TabContent) []*TabClose {
	var tabCloses []*TabClose
	for _, tabContent := range tabContents {
		tabClose := NewTabClose(tabContent.Title, tabContent.Tooltip+" "+tabContent.Title, tabContent.Closeable, tabContent.Panel)
		tabCloses = append(tabCloses, tabClose)
	}
	return tabCloses
}

func NewTabClose(title string, tooltip string, closeable bool, panel unison.Paneler) *TabClose {
	d := &TabClose{
		Panel:       unison.Panel{},
		title:       title,
		tooltip:     tooltip,
		closeable:   closeable,
		ScrollPanel: nil,
	}

	d.Self = d
	SetScrollLayout(d, 1)
	d.AddChild(panel)
	return d
}

func (t *TabClose) SetColumnHeader(header unison.Paneler) {
	mylog.CheckNil(t.ScrollPanel)
	t.ScrollPanel.SetColumnHeader(header)
}

func (t *TabClose) MayAttemptClose() bool {
	return t.closeable
}

func (t *TabClose) AttemptClose() bool {
	if dc := unison.Ancestor[*unison.DockContainer](t); dc != nil {
		dc.Close(t)
		return true
	}
	return false
}

func (t *TabClose) Modified() bool {
	return false // todo update title
}

func (t *TabClose) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  unison.DocumentSVG,
		Size: suggestedSize,
	}
}

func (t *TabClose) Title() string {
	return i18n.Text(t.title)
}

func (t *TabClose) Tooltip() string {
	return i18n.Text(t.tooltip)
}

var _ unison.Dockable = &Tab{}

type Tab struct {
	unison.Panel
	title     string
	tooltip   string
	closeable bool
}

func NewTabWithTable[T any](table *Node[T], header *TableHeader[T], title string, tooltip string, closeable bool) *Tab {
	panel := NewPanel()
	panel.AddChild(table)
	panel.AddChild(header)

	scrollPanelFill := NewScrollPanelFill(panel)
	scrollPanelFill.SetColumnHeader(header)

	tab := NewTab(title, tooltip, closeable, scrollPanelFill)
	return tab
}

func NewTab(title string, tooltip string, closeable bool, panel unison.Paneler) *Tab {
	d := &Tab{
		Panel:     unison.Panel{},
		title:     title,
		tooltip:   tooltip,
		closeable: closeable,
	}
	d.Self = d
	SetScrollLayout(d, 1)
	d.AddChild(panel)
	return d
}

func NewTabs(tabContents ...TabContent) []*Tab { // todo delete
	var tabs []*Tab
	for _, tabContent := range tabContents {
		tab := NewTab(tabContent.Title, tabContent.Tooltip+" "+tabContent.Title, tabContent.Closeable, tabContent.Panel)
		tabs = append(tabs, tab)
	}
	return tabs
}

func (t *Tab) SetTitle(title string) {
	t.title = title
}

func (t *Tab) SetTooltip(tooltip string) {
	t.tooltip = tooltip
}

func (t *Tab) SetCloseable(closeable bool) {
	t.closeable = closeable
}

func (t *Tab) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  unison.DocumentSVG,
		Size: suggestedSize,
	}
}

func (t *Tab) Title() string {
	return i18n.Text(t.title)
}

func (t *Tab) Tooltip() string {
	return i18n.Text(t.tooltip)
}

func (t *Tab) Modified() bool {
	return false // todo update title
}
