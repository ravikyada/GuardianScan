package handlers

import (
	"Nova/scanner"
	"Nova/storage"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LogsHandler streams logs for live updates during the scan and saves the result to reports.json
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Extract the target URL from the query parameters
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		fmt.Fprintf(w, "data: ❌ No target URL provided.\n\n")
		w.(http.Flusher).Flush()
		return
	}

	// Check if the site is reachable
	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Fprintf(w, "data: ❌ Site is not reachable: %v\n\n", err)
		w.(http.Flusher).Flush()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "data: ❌ Site returned status code %d\n\n", resp.StatusCode)
		w.(http.Flusher).Flush()
		return
	}

	// Parse hostname
	hostname := strings.TrimPrefix(targetURL, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.Split(hostname, "/")[0]

	// Store logs
	var collectedLogs []string

	// Stream logs
	addLog := func(log string) {
		collectedLogs = append(collectedLogs, log) // Collect log
		fmt.Fprintf(w, "data: %s\n\n", log)        // Stream log to client
		w.(http.Flusher).Flush()
	}

	addLog(fmt.Sprintf("🔍 Starting scan for %s...", targetURL))
	time.Sleep(500 * time.Millisecond)

	// Scan open ports
	addLog("🔍 Checking open ports...")
	portResults := scanner.ScanPorts(hostname)
	for port, isOpen := range portResults {
		if isOpen {
			addLog(fmt.Sprintf("❌ Port %d is open", port))
		} else {
			addLog(fmt.Sprintf("✅ Port %d is closed", port))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Scan HTTP security headers
	addLog("🔍 Checking security headers...")
	headerResults := scanner.CheckSecurityHeaders(targetURL)
	for header, exists := range headerResults {
		if exists {
			addLog(fmt.Sprintf("✅ Security header found: %s", header))
		} else {
			addLog(fmt.Sprintf("❌ Missing security header: %s", header))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Check for sensitive files
	// Check for sensitive files
	addLog("🔍 Checking sensitive files...")
	fileResults := scanner.CheckSensitiveFiles(targetURL)
	for category, files := range fileResults {
		addLog(fmt.Sprintf("📁 Checking %s:", category))
		for file, exists := range files {
			if exists == "accessible" {
				fullURL := fmt.Sprintf("%s/%s", targetURL, file)
				addLog(fmt.Sprintf("  ❌ Found: %s (%s)", file, fullURL))
			} else {
				addLog(fmt.Sprintf("  ✅ Not accessible: %s", file))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Final completion log
	addLog("🎉 Scan completed successfully!")

	// Save logs to reports.json
	report := storage.ScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		URL:       targetURL,
		Results:   strings.Join(collectedLogs, "\n"),
	}
	storage.SaveReport(report)
}
