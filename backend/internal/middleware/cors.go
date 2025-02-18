package middleware

import (
	"net/http"
)

// SimpleCORS sets basic headers allowing any origin to make requests.
// For local dev, "*" is fine, but in production you might restrict it.
func SimpleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's an OPTIONS request (the "preflight" check), respond OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, continue to the actual handler
		next.ServeHTTP(w, r)
	})
}
