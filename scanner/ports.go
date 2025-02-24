package scanner

import (
	"fmt"
	"net"
	"time"
)

var commonPorts = map[int]string{
	22:   "SSH",
	3306: "MySQL",
	5432: "PostgreSQL",
	1433: "MSSQL",
	6379: "Redis",
	80:   "HTTP",
	443:  "HTTPS",
}

// ScanPorts checks for open ports on the target host
func ScanPorts(host string) map[int]bool {
	results := make(map[int]bool)
	for port := range commonPorts {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			results[port] = true // Port is open
			conn.Close()
		} else {
			results[port] = false // Port is closed or filtered
		}
	}
	return results
}
