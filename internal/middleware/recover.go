package middleware

import "net/http"

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				//TODO: log
			}
		}()
		next.ServeHTTP(w, r)
	})
}
