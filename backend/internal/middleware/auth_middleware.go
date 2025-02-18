package middleware

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// WithFirebaseAuth is a middleware that ensures requests have a valid Firebase ID token.
func WithFirebaseAuth(firebaseAuth *auth.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := firebaseAuth.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If token is valid, proceed with the original handler
		next.ServeHTTP(w, r)
	}
}
