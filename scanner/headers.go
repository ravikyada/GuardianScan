package scanner

import (
	"net/http"
)

var securityHeaders = map[string]string{
	"Strict-Transport-Security": "HSTS protection",
	"X-Frame-Options":           "Prevents Clickjacking",
	"X-Content-Type-Options":    "Prevents MIME-type sniffing",
	"Content-Security-Policy":   "Prevents XSS attacks",
	"Referrer-Policy":           "Controls referrer data exposure",
	"Permissions-Policy":        "Restricts APIs based on origin",
}

// CheckSecurityHeaders scans for security headers in the target URL
func CheckSecurityHeaders(url string) map[string]bool {
	results := make(map[string]bool)
	resp, err := http.Get(url)
	if err != nil {
		return results
	}
	defer resp.Body.Close()

	for header := range securityHeaders {
		if _, exists := resp.Header[header]; exists {
			results[header] = true
		} else {
			results[header] = false
		}
	}
	return results
}
