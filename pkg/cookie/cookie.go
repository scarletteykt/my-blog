package cookie

import (
	"errors"
	"net/http"
	"strings"
)

const IDCookieName = "idCookie"
const IDCookieSep = ":"
const IDCookieMaxAge = 300_000_000

type IDCookie struct {
	Username string
	Sign     string
	Cookie   *http.Cookie
}

func NewIDCookie(username string, sign string) *IDCookie {
	value := username + IDCookieSep + sign

	cookie := &http.Cookie{
		Name:   IDCookieName,
		Value:  value,
		MaxAge: IDCookieMaxAge,
	}

	return &IDCookie{
		Username: username,
		Sign:     sign,
		Cookie:   cookie,
	}
}

func ParseFromCookie(cookie *http.Cookie) (*IDCookie, error) {
	value := cookie.Value

	if value == "" {
		return nil, errors.New("empty cookie value")
	}

	parsed := strings.Split(value, IDCookieSep)

	if len(parsed) != 2 {
		return nil, errors.New("invalid cookie value")
	}

	username := parsed[0]
	sign := parsed[1]

	return &IDCookie{
		Username: username,
		Sign:     sign,
		Cookie:   cookie,
	}, nil
}
