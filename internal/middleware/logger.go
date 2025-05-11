package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"метод":    r.Method,
			"путь":     r.URL.Path,
			"Значение": time.Since(start).String(),
		}).Info("Обработанный запрос")
	})
}
