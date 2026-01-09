package middleware

import (
	"net/http"
)

// CORSMiddleware creates HTTP middleware for CORS handling
func CORSMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	if allowedOrigins == "" {
		allowedOrigins = "*" // Default for development
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Continue with next handler
			next.ServeHTTP(w, r)
		})
	}
}
