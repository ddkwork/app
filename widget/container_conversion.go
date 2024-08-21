package widget

import (
	"strings"

	"github.com/ddkwork/unison"
)

func InstallContainerConversionHandlers[T any](paneler unison.Paneler, table *Node[T]) {
	p := paneler.AsPanel()
	p.InstallCmdHandlers(0,
		func(_ any) bool { return CanConvertToContainer(table) },
		func(_ any) { ConvertToContainer(table) })
	p.InstallCmdHandlers(0,
		func(_ any) bool { return CanConvertToNonContainer(table) },
		func(_ any) { ConvertToNonContainer(table) })
}

func CanConvertToContainer[T any](table *Node[T]) bool {
	for _, row := range table.SelectedRows(false) {
		if data := row; data != nil && !data.Container() {
			return true
		}
	}
	return false
}

func CanConvertToNonContainer[T any](table *Node[T]) bool {
	for _, row := range table.SelectedRows(false) {
		if data := row; data != nil && data.Container() && !data.HasChildren() {
			return true
		}
	}
	return false
}

func ConvertToContainer[T any](table *Node[T]) {
	before := &Node[T]{}
	after := &Node[T]{}
	for _, row := range table.SelectedRows(false) {
		if data := row; data != nil && !data.Container() {
			data.SetType(data.GetType() + ContainerKeyPostfix)
			data.SetOpen(true)
		}
	}
	if len(before.Children) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			mgr.Add(&unison.UndoEdit[*Node[T]]{
				ID:         unison.NextUndoID(),
				EditName:   "ConvertToContainer",
				UndoFunc:   func(edit *unison.UndoEdit[*Node[T]]) { edit.BeforeData.SetType(edit.BeforeData.GetType()) },
				RedoFunc:   func(edit *unison.UndoEdit[*Node[T]]) { edit.AfterData.SetType(edit.AfterData.GetType()) },
				BeforeData: before,
				AfterData:  after,
			})
		}
		table.SetRootRows(table.Children)
	}
}

func ConvertToNonContainer[T any](table *Node[T]) {
	before := &Node[T]{}
	after := &Node[T]{}
	for _, row := range table.SelectedRows(false) {
		if data := row; data != nil && data.Container() && !data.HasChildren() {
			data.SetType(strings.TrimSuffix(data.GetType(), ContainerKeyPostfix))
		}
	}
	if len(before.Children) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			mgr.Add(&unison.UndoEdit[*Node[T]]{
				ID:         unison.NextUndoID(),
				EditName:   "ConvertToNonContainer",
				UndoFunc:   func(edit *unison.UndoEdit[*Node[T]]) { edit.BeforeData.SetType(edit.BeforeData.GetType()) },
				RedoFunc:   func(edit *unison.UndoEdit[*Node[T]]) { edit.AfterData.SetType(edit.AfterData.GetType()) },
				BeforeData: before,
				AfterData:  after,
			})
		}
		table.SetRootRows(table.Children)
	}
}
