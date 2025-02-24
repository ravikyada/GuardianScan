package main

import (
	"Nova/handlers"
	"log"
	"net/http"
)

func main() {
	log.Println("📢 Starting Security Scanner...")
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/scan", handlers.ScanHandler)
	http.HandleFunc("/reports", handlers.ReportsHandler)
	http.HandleFunc("/logs", handlers.LogsHandler)

	log.Println("🚀 Security Scanner running on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
