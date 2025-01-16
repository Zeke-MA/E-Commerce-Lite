package middleware

import "net/http"

func LogIncomingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to open the log file
		// Generate the log file structure
		// Capture the fields if available
		// Write to the log
		// Save the log
		// defer close the file
		// serve the http endpoint
	})
}
