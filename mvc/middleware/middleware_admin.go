package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func IsAdminLoggedIn(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Admin"] != nil {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Admin girişi yapılmadı")
		}
	})
}

func MiddlewareForAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("admin") != "" {
			r.Header.Set("Admin", os.Getenv("admin"))
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middleware admin")
	})
}
