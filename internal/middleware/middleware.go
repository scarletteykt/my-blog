package middleware

import (
	mw "github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Middleware() []func(handler http.Handler) http.Handler {
	mws := []func(handler http.Handler) http.Handler{
		mw.RequestID,
		mw.RealIP,
		mw.Logger,
		mw.Recoverer,
	}

	return mws
}
