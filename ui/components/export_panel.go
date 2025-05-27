// ui/components/export_panel.go - GUI export control panel
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ExportPanel struct {
	window fyne.Window
}

func NewExportPanel(window fyne.Window) *ExportPanel {
	return &ExportPanel{window: window}
}

func (e *ExportPanel) SetOnCSVExport(callback func(string)) {
	// Implementation for CSV export button
}

func (e *ExportPanel) SetOnJSONExport(callback func(string)) {
	// Implementation for JSON export button
}

func (e *ExportPanel) Render() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewButton("Export CSV", func() {
			dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
				if uri != nil {
					// Trigger CSV export
				}
			}, e.window)
		}),
		widget.NewButton("Export JSON", func() {
			dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
				if uri != nil {
					// Trigger JSON export
				}
			}, e.window)
		}),
	)
}
