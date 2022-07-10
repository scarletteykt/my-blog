package middleware

import (
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"net/http"
)

func Middleware(opts ...Option) []func(handler http.Handler) http.Handler {
	o := initOptions(opts)
	c := cors.New(o.corsOptions)
	a := auth.New(o.authOptions)

	mws := []func(handler http.Handler) http.Handler{
		c.Handler,
		a.Handler,
		mw.RequestID,
		mw.RealIP,
		mw.Logger,
		Recover,
	}

	return mws
}
