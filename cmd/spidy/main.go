package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/twiny/spidy/tree/main/cmd/spidy/api"
	"github.com/twiny/spidy/tree/main/internal/pkg/spider/v1"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ConfigRequest struct {
	ConfigPath string   `json:"configPath"`
	URLs       []string `json:"urls"`
}

type Result struct {
	Domain string `json:"domain"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func main() {
	// Serve static files (React app)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API to start scraping
	http.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		spider, err := api.NewSpider(req.ConfigPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go spider.Shutdown()

		go func() {
			if err := spider.Start(req.URLs); err != nil {
				log.Println(err)
			}
		}()

		w.WriteHeader(http.StatusOK)
	})

	// WebSocket for streaming results
	http.HandleFunc("/api/results", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		// Register connection with spider
		api.RegisterWebSocket(conn)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
