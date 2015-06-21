package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logrus.Debugf("%s - %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
}
