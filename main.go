package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Structure to hold IP and path combination count
type requestCounter struct {
	mu       sync.Mutex
	requests map[string]map[string]int
}

func newRequestCounter() *requestCounter {
	return &requestCounter{
		requests: make(map[string]map[string]int),
	}
}

// Count the IP and path combination
func (rc *requestCounter) count(ip, path string) int {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// If IP entry does not exist, create it
	if _, ok := rc.requests[ip]; !ok {
		rc.requests[ip] = make(map[string]int)
	}

	// Increment the counter for the given path
	rc.requests[ip][path]++
	return rc.requests[ip][path]
}

// Get all IP and path combinations with counts
func (rc *requestCounter) getAllCounts() map[string]map[string]int {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Make a deep copy to avoid race conditions
	copy := make(map[string]map[string]int)
	for ip, paths := range rc.requests {
		copy[ip] = make(map[string]int)
		for path, count := range paths {
			copy[ip][path] = count
		}
	}
	return copy
}

// Helper function to get the client's IP address
func getClientIP(r *http.Request) string {
	// Check if X-Forwarded-For header is present
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, we take the first one
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	// If X-Forwarded-For is not present, use RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

func main() {
	counter := newRequestCounter()

	// Handler function for all paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the client IP, considering possible proxy headers
		ip := getClientIP(r)
		path := r.URL.Path
		method := r.Method
		forwardedFor := r.Header.Get("X-Forwarded-For")
		queryString := r.URL.RawQuery

		// Retrieve the last hop's IP (direct connection to the server)
		lasthop, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			lasthop = ""
		}

		// Retrieve cookie names (not logged to StdOut but included in JSON response)
		var cookieNames []string
		for _, cookie := range r.Cookies() {
			cookieNames = append(cookieNames, cookie.Name)
		}

		// If "/listall" is requested, return all data
		if path == "/listall" {
			allCounts := counter.getAllCounts()
			jsonResponse(w, allCounts)
			return
		}

		// Count and get the current count for the specific IP and path
		count := counter.count(ip, path)

		// Get the current time in ISO-8601 format
		currentTime := time.Now().Format(time.RFC3339)

		// Log the requested information to standard output
		fmt.Printf("%s - %s - %s - %s - %s - %d - %s\n", currentTime, lasthop, forwardedFor, method, path, count, queryString)

		// Prepare the JSON response
		response := map[string]interface{}{
			"ip":             ip,
			"path":           path,
			"method":         method,
			"count":          count,
			"x_forwarded_for": forwardedFor,
			"query_string":   queryString,
			"cookie_names":   cookieNames,
			"lasthop":        lasthop,
		}

		jsonResponse(w, response)
	})

	fmt.Println("Starting server on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}

// Helper function to respond with JSON
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

