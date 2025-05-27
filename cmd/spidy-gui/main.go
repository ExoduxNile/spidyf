// cmd/spidy-gui/main.go - Main entry point for GUI application
package main

import (
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/twiny/spidy/v2/ui/components"
)

func main() {
	// Create GUI app
	myApp := app.New()
	window := myApp.NewWindow("Spidy Domain Scraper")
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

	// Start health check server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080" // Default port for local development
		}

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			// Optional: log or handle error
			panic("Failed to start health server: " + err.Error())
		}
	}()

	// Run GUI (blocking)
	window.ShowAndRun()
}

