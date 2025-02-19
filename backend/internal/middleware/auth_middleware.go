package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type UIDCtxKey struct{}

// WithFirebaseAuth is a middleware that ensures requests have a valid Firebase access token.
// The middleware also updates the context with the user ID (UID) retrieved from the verified token.
func WithFirebaseAuth(firebaseAuth *auth.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		ctx := r.Context()

		token, err := firebaseAuth.VerifyIDToken(ctx, tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UIDCtxKey{}, token.UID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
