package v1

import (
	"encoding/json"
	"errors"
	"github.com/scraletteykt/my-blog/internal/user"
	"github.com/scraletteykt/my-blog/pkg/bcrypt"
	"github.com/scraletteykt/my-blog/pkg/cookie"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"github.com/scraletteykt/my-blog/pkg/server"
	"github.com/scraletteykt/my-blog/pkg/sign"
	"net/http"
)

type signUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) SignUp(w http.ResponseWriter, r *http.Request) {
	var s signUpInput

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		logger.Warnf("user sign up: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.Hash(s.Password)
	if err != nil {
		logger.Errorf("user sign up: bcrypt hashing error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}

	_, err = a.users.CreateUser(user.CreateUser{
		Username:     s.Username,
		PasswordHash: hashed,
	})
	if err == user.ErrUserAlreadyExists {
		logger.Warnf("user sign up: user already exist error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		logger.Errorf("user sign up: create user error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}

	server.ResponseJSONWithCode(w, r, http.StatusOK, nil)
}

func (a *API) SignIn(w http.ResponseWriter, r *http.Request) {
	var s signInInput

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		logger.Warnf("user sign in: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}

	u, err := a.users.GetUser(s.Username)
	if err == user.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusUnauthorized, err)
		return
	}
	if err != nil {
		logger.Errorf("user sign in: get user error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}

	err = bcrypt.Compare(u.PasswordHash, s.Password)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusUnauthorized, errors.New("wrong username or password"))
		return
	}

	signer := sign.NewSigner("deadbeef")
	http.SetCookie(w, cookie.NewIDCookie(u.Username, signer.EncodeBase64(signer.Sign(u.Username))).Cookie)

	server.ResponseJSONWithCode(w, r, http.StatusOK, nil)
}
