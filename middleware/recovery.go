package middleware

import (
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(
					"Восстановлена паника при обработке %s %s: %v",
					r.Method,
					r.URL.Path,
					err,
				)

				http.Error(
					w,
					"Внутренняя ошибка сервера",
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
