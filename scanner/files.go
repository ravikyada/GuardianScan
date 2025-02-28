package scanner

import (
	"fmt"
	"net/http"
	"time"
)

// Define categories of sensitive files
var databaseFiles = []string{
	".db",        // Generic database file
	".sql",       // SQL database backup
	".sqlite",    // SQLite database file
	".mongodb",   // MongoDB DB file
	"backup.sql", // SQL backup file
	"dump.sql",   // SQL dump file
}

var wordpressFiles = []string{
	".htaccess",     // WordPress configuration
	"wp-config.php", // WordPress configuration
	".git/config",   // Git config file
	".gitignore",    // Git ignore file
}

var nodejsFiles = []string{
	"package.json", // Node.js package configuration
	".env",         // Environment variable file for Node.js
	".env.local",   // Local environment variables for Node.js
	"config.json",  // JSON config file for Node.js apps
}

var vulnerabilityFiles = []string{
	".env",               // Environment file
	"config.php",         // PHP config file
	"config.json",        // JSON config file
	"phpinfo.php",        // PHP info (Exposes sensitive server info)
	"access.log",         // Access logs
	"error.log",          // Error logs
	"secret.key",         // Rails secret key file
	"config/secrets.yml", // Rails secrets file
	"database.yml",
	"web/js/vans.js", // Database config (Rails or other frameworks)
}

// CheckSensitiveFiles scans for publicly accessible sensitive files and categorizes them
func CheckSensitiveFiles(url string) map[string]map[string]string {
	// A map to hold the categorized files and their status
	results := make(map[string]map[string]string)

	// Initialize the categories in the results map
	results["Database Files"] = make(map[string]string)
	results["WordPress Files"] = make(map[string]string)
	results["Node.js Files"] = make(map[string]string)
	results["Common Vulnerability Files"] = make(map[string]string)

	client := http.Client{
		Timeout: 5 * time.Second, // Set a timeout for each HTTP request
	}

	// Helper function to categorize and check files
	checkFiles := func(fileList []string, category string) {
		for _, file := range fileList {
			fileURL := fmt.Sprintf("%s/%s", url, file)
			resp, err := client.Head(fileURL)
			if err != nil {
				// Log the error if the request fails
				fmt.Printf("Error accessing %s: %v\n", fileURL, err)
				results[category][file] = "Secure (Error)" // Mark file as secure if there's an error
			} else {
				// Check the response status code and append it to the results
				statusText := fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				if resp.StatusCode == http.StatusOK {
					results[category][file] = "Accessible (" + statusText + ")" // File is publicly accessible
				} else {
					results[category][file] = "Secure (" + statusText + ")" // File is secure if not accessible
				}
				resp.Body.Close() // Close the response body
			}
		}
	}

	// Check each category of files
	checkFiles(databaseFiles, "Database Files")
	checkFiles(wordpressFiles, "WordPress Files")
	checkFiles(nodejsFiles, "Node.js Files")
	checkFiles(vulnerabilityFiles, "Common Vulnerability Files")

	return results
}
