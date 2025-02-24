package handlers

import (
	"Nova/storage"
	"encoding/json"
	"html/template"
	"net/http"
)

// ReportsHandler serves the reports page and JSON data
func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the reports HTML page
		tmpl := template.Must(template.ParseFiles("templates/reports.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Return the JSON data for reports
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storage.LoadReports())
	}
}
