package widget

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/maps"
	"github.com/richardwilkes/unison"
)

type idMgr struct {
	index int
	idMao *maps.SafeMap[int, bool]
}

func (i *idMgr) add() int {
	i.index++
	if i.idMao.Has(i.index) {
		mylog.Check("detected duplicate id")
	}
	i.idMao.Set(i.index, true)
	//	mylog.Trace("add action id", i.index)
	return i.index
}

func newIdMgr() *idMgr {
	return &idMgr{
		index: unison.UserBaseID,
		idMao: new(maps.SafeMap[int, bool]),
	}
}

var defaultIdMgr = newIdMgr()

type (
	ContextMenuItem struct {
		Title string
		id    int
		Can   func(any) bool
		Do    func(any)
	}
	MouseDownCallbackType func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool
	ContextMenuItems      struct {
		Panel            *unison.Panel
		DefaultMouseDown MouseDownCallbackType
		items            []ContextMenuItem
	}
)

func NewContextMenuItems(panel unison.Paneler, DefaultMouseDown MouseDownCallbackType, items ...ContextMenuItem) *ContextMenuItems {
	for i := 0; i < len(items); i++ {
		if items[i].Can == nil {
			items[i].Can = unison.AlwaysEnabled
		}
		items[i].id = defaultIdMgr.add()
	}
	label, ok := panel.(*unison.Label)
	if ok {
		LabelStyle(label)
	}
	if DefaultMouseDown == nil {
		mylog.CheckNil(panel.AsPanel().MouseDownCallback)
		DefaultMouseDown = panel.AsPanel().MouseDownCallback
	}
	return &ContextMenuItems{
		Panel:            panel.AsPanel(),
		DefaultMouseDown: DefaultMouseDown,
		items:            items,
	}
}

func InsertCmdContextMenuItem(panel *unison.Panel, title string, cmdID int, id *int, cm unison.Menu) {
	if panel.CanPerformCmd(panel, cmdID) {
		useID := *id
		*id++
		cm.InsertItem(-1, cm.Factory().NewItem(unison.PopupMenuTemporaryBaseID+useID, title, unison.KeyBinding{}, nil,
			func(_ unison.MenuItem) {
				panel.PerformCmd(panel, cmdID)
			}))
	}
}

func (c *ContextMenuItems) Install() {
	defer func() {
		for _, item := range c.items {
			c.Panel.InstallCmdHandlers(item.id, item.Can, item.Do)
		}
	}()
	c.Panel.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		stop := c.DefaultMouseDown(where, button, clickCount, mod)
		if button == unison.ButtonRight && clickCount == 1 && !c.Panel.Window().InDrag() {
			f := unison.DefaultMenuFactory()
			cm := f.NewMenu(unison.PopupMenuTemporaryBaseID|unison.ContextMenuIDFlag, "", nil)
			id := 1
			for _, one := range c.items {
				if one.id == -1 {
					cm.InsertSeparator(-1, true)
				} else {
					InsertCmdContextMenuItem(c.Panel, one.Title, one.id, &id, cm)
				}
			}
			count := cm.Count()
			if count > 0 {
				count--
				if cm.ItemAtIndex(count).IsSeparator() {
					cm.RemoveItem(count)
				}
				c.Panel.FlushDrawing()
				cm.Popup(unison.Rect{
					Point: c.Panel.PointToRoot(where),
					Size: unison.Size{
						Width:  1,
						Height: 1,
					},
				}, 0)
			}
			cm.Dispose()
		}
		return stop
	}
}
