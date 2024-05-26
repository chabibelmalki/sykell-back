package middleware

import (
	"fmt"
	"net/http"
)

// APIKey - To modify
const APIKey = "STATICAPIkeyForTestingAnAuthMethod"

// AuthMiddleware is a middleware function that checks for a valid API key in the request header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the API key from the Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the API key is present and matches the expected value
		if authHeader != fmt.Sprintf("Bearer %s", APIKey) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: Invalid API key"))
			return
		}

		// If the key is valid, call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
