package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"Nova/storage"
)

// IndexHandler serves the single-page application (index.html)
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Parse and execute the index.html template instead of a separate reports.html
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// FetchReportsHandler returns JSON data for generated reports
func FetchReportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set JSON response header
	w.Header().Set("Content-Type", "application/json")

	// Load reports from storage (generated reports)
	reports := storage.LoadReports()

	// Encode and send reports as JSON
	json.NewEncoder(w).Encode(reports)
}
