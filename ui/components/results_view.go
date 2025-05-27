// ui/components/results_view.go - Domain results display component
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ResultsView struct {
	data  binding.UntypedList
	table *widget.Table
}

func NewResultsView() *ResultsView {
	view := &ResultsView{
		data: binding.NewUntypedList(),
	}
	
	view.table = widget.NewTable(
		func() (int, int) {
			length, _ := view.data.Length()
			return length, 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// Cell rendering logic
		},
	)
	
	return view
}

func (v *ResultsView) AddDomain(domain interface{}) {
	v.data.Append(domain)
	v.table.Refresh()
}

func (v *ResultsView) Render() fyne.CanvasObject {
	return container.NewScroll(v.table)
}
