package middleware

import "net/http"

func Header(values map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range values {
				w.Header().Set(key, value)
			}
			next.ServeHTTP(w, r)
		})
	}
}
