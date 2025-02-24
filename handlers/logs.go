package handlers

import (
	"Nova/scanner"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// checkURLStatus validates if the given URL is accessible
func checkURLStatus(targetURL string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// extractHostname extracts the hostname from a given URL
func extractHostname(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil // Extract only the domain/hostname
}

// LogsHandler streams logs for live updates during the scan
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Ensure Flusher is available
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Extract target URL
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		fmt.Fprintf(w, "data: ❌ No target URL provided.\n\n")
		flusher.Flush()
		return
	}

	// Check if URL is reachable
	if !checkURLStatus(targetURL) {
		fmt.Fprintf(w, "data: ❌ The provided URL is not accessible (Non-200 response).\n\n")
		flusher.Flush()
		return
	}

	// Extract hostname for scanning
	hostname, err := extractHostname(targetURL)
	if err != nil {
		fmt.Fprintf(w, "data: ❌ Failed to extract hostname from URL.\n\n")
		flusher.Flush()
		return
	}

	// Store logs
	var collectedLogs []string
	addLog := func(log string) {
		collectedLogs = append(collectedLogs, log)
		fmt.Fprintf(w, "data: %s\n\n", log)
		flusher.Flush()
	}

	// Start scan
	addLog(fmt.Sprintf("🔍 Starting scan for %s...", targetURL))
	time.Sleep(500 * time.Millisecond)

	// Check open ports
	addLog("🔍 Checking open ports...")
	portResults := scanner.ScanPorts(hostname) // Use extracted hostname instead of full URL
	for port, isOpen := range portResults {
		if isOpen {
			addLog(fmt.Sprintf("❌ Port %d is open", port))
		} else {
			addLog(fmt.Sprintf("✅ Port %d is closed", port))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Check security headers
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

	// Check sensitive files
	addLog("🔍 Checking sensitive files...")
	fileResults := scanner.CheckSensitiveFiles(targetURL)
	for file, accessible := range fileResults {
		if accessible {
			addLog(fmt.Sprintf("❌ Sensitive file accessible: %s", file))
		} else {
			addLog(fmt.Sprintf("✅ Sensitive file secure: %s", file))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Final completion log
	addLog("🎉 Scan completed successfully!")

	// Explicitly signal the client that the connection will close
	addLog("🔴 Scan complete, closing connection.")

	// Close SSE connection gracefully
	fmt.Fprintf(w, "data: 🔴 SSE Connection Closed.\n\n")
	flusher.Flush()
	time.Sleep(1 * time.Second)
}
