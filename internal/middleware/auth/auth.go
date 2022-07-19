package auth

import (
	"github.com/scraletteykt/my-blog/internal/service"
	"github.com/scraletteykt/my-blog/pkg/auth"
	"github.com/scraletteykt/my-blog/pkg/cookie"
	signer "github.com/scraletteykt/my-blog/pkg/sign"
	"net/http"
)

type Config struct {
	Secret string
}

type Auth struct {
	secret string
	users  *service.UsersService
}

func New(cfg *Config, users *service.UsersService) *Auth {
	return &Auth{
		secret: cfg.Secret,
		users:  users,
	}
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := signer.NewSigner(a.secret)
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

		u, err := a.users.GetUser(r.Context(), idCookie.Username)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(auth.WithUser(r.Context(), auth.User{
			ID:       u.ID,
			Username: u.Username,
		})))
	})
}
