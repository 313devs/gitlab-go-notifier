package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

// Auth is a middleware that checks if the request has a valid token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth middleware")
		token := r.Header.Get("Authorization")
		if token != os.Getenv("TOKEN") || token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}
