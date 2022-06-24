package storage

import "strings"

// CheckEndpoint check  whether endpoint is combine with ip:port
// example: http://127.0.0.1:9999/xxx
func CheckEndpoint(endpoint string) bool {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	return strings.ContainsAny(endpoint, ":")
}
