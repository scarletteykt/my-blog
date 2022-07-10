package auth

import (
	"github.com/scraletteykt/my-blog/internal/user"
	"github.com/scraletteykt/my-blog/pkg/cookie"
	signer "github.com/scraletteykt/my-blog/pkg/sign"
	"net/http"
)

type Auth struct {
	Options Options
}

func New(opts Options) *Auth {
	return &Auth{
		Options: opts,
	}
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := signer.NewSigner(a.Options.Secret)
		httpCookie, err := r.Cookie(cookie.IDCookieName)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		idCookie, err := cookie.ParseFromCookie(httpCookie)
		sign, decodeErr := s.DecodeBase64(idCookie.Sign)

		if err != nil || decodeErr != nil || !s.Verify(sign, idCookie.Username) {
			next.ServeHTTP(w, r)
			return
		}

		u, err := a.Options.Services.Users.GetUser(idCookie.Username)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(user.WithUser(r.Context(), u)))
	})
}
