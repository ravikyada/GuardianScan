package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type ScanResult struct {
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Results   string `json:"results"`
}

var (
	reportFile  = "reports.json"
	reportMutex = sync.Mutex{}
)

func LoadReports() []ScanResult {
	var reports []ScanResult

	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		return reports // Return empty if the file doesn't exist
	}

	data, err := ioutil.ReadFile(reportFile)
	if err != nil {
		return reports // Return empty on error
	}

	json.Unmarshal(data, &reports)
	return reports
}

// func SaveReport(newReport ScanResult) {
// 	reportMutex.Lock()
// 	defer reportMutex.Unlock()

// 	reports := LoadReports()
// 	reports = append(reports, newReport)

//		data, _ := json.MarshalIndent(reports, "", "  ")
//		ioutil.WriteFile(reportFile, data, 0644)
//	}
func SaveReport(newReport ScanResult) {
	reportMutex.Lock()
	defer reportMutex.Unlock()

	// Load existing reports
	reports := LoadReports()

	// Append the new report
	reports = append(reports, newReport)

	// Save updated reports back to the file
	data, _ := json.MarshalIndent(reports, "", "  ")
	err := ioutil.WriteFile(reportFile, data, 0644)
	if err != nil {
		fmt.Printf("Error saving report: %v\n", err) // Add debug logs
	}
}
