// cmd/spidy-gui/main.go - Main entry point for GUI application
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/twiny/spidy/v2/ui/components"
)

func main() {
	app := app.New()
	window := app.NewWindow("Spidy Domain Scraper")
	window.Resize(fyne.NewSize(1200, 800))

	// Initialize components
	exportPanel := components.NewExportPanel()
	resultsView := components.NewResultsView()

	// Set up main layout
	content := container.NewBorder(
		nil,
		exportPanel.Render(),
		nil,
		nil,
		resultsView.Render(),
	)

	window.SetContent(content)
	window.ShowAndRun()
}
