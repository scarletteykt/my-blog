package server

import (
	"github.com/go-chi/render"
	"net/http"
)

type httpError struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
	DebugID string        `json:"debug_id,omitempty"`
}

type ErrorDetail struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

func ResponseJSON(w http.ResponseWriter, r *http.Request, obj interface{}) {
	if obj == nil {
		obj = struct {
		}{}
	}
	render.JSON(w, r, obj)
}

func ResponseJSONWithCode(w http.ResponseWriter, r *http.Request, code int, obj interface{}) {
	if obj == nil {
		obj = struct {
		}{}
	}
	render.Status(r, code)
	render.JSON(w, r, obj)
}

func ErrorJSON(w http.ResponseWriter, r *http.Request, code int, err error, errs ...ErrorDetail) {
	var resp httpError

	resp.Meta.Code = code
	resp.Meta.Message = err.Error()
	resp.Meta.Details = errs

	render.Status(r, code)
	render.JSON(w, r, &resp)
}
