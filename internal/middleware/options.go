package middleware

import (
	"github.com/go-chi/cors"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"net/http"
)

type options struct {
	corsOptions      cors.Options
	authOptions      auth.Options
	customMiddleware []func(handler http.Handler) http.Handler
}

type Option func(*options)

func initOptions(opts []Option) *options {
	o := &options{}
	for i := range opts {
		opts[i](o)
	}
	return o
}

func WithCORSOptions(corsOpts cors.Options) Option {
	return func(opts *options) {
		opts.corsOptions = corsOpts
	}
}

func WithAuthOptions(authOpts auth.Options) Option {
	return func(opts *options) {
		opts.authOptions = authOpts
	}
}
