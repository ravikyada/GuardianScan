package scanner

import (
	"fmt"
	"net/http"
)

var sensitiveFiles = []string{
	".htaccess",
	".env",
	".git/config",
	"backup.sql",
	"wp-config.php",
}

// CheckSensitiveFiles scans for publicly accessible sensitive files
func CheckSensitiveFiles(url string) map[string]bool {
	results := make(map[string]bool)
	for _, file := range sensitiveFiles {
		fileURL := fmt.Sprintf("%s/%s", url, file)
		resp, err := http.Head(fileURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			results[file] = true // File is publicly accessible
		} else {
			results[file] = false // File is secure
		}
	}
	return results
}
