package handlers

import (
	"Nova/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	targetURL := r.FormValue("url")

	// Perform the scan (dummy data for testing)
	results := map[string]interface{}{
		"open_ports":       map[int]bool{80: true, 443: false}, // Example
		"security_headers": map[string]bool{"Strict-Transport-Security": true},
		"sensitive_files":  map[string]bool{".env": false},
	}

	// Save the report
	report := storage.ScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		URL:       targetURL,
		Results:   string(toJSON(results)),
	}
	fmt.Printf("Saving report: %+v\n", report) // Debug log
	storage.SaveReport(report)

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func toJSON(data interface{}) string {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	return string(bytes)
}
