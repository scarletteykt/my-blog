package v1

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/internal/service"
	"github.com/scraletteykt/my-blog/pkg/server"
	"net/http"
	"strconv"
)

type createTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type updateTag struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

func (a *API) GetTagByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "tagID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	t, err := a.tags.GetTagByID(r.Context(), int(id))
	if err == service.ErrNotFound {
		server.ResponseJSONWithCode(w, r, http.StatusNoContent, struct{}{})
		return
	}
	if err != nil {
		a.log.Errorf("error: get tag by id: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, t)
}

func (a *API) GetTags(w http.ResponseWriter, r *http.Request) {
	t, err := a.tags.GetTags(r.Context())
	if err == service.ErrNotFound {
		server.ResponseJSONWithCode(w, r, http.StatusNoContent, struct{}{})
		return
	}
	if err != nil {
		a.log.Errorf("error: get tags: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, t)
}

func (a *API) CreateTag(w http.ResponseWriter, r *http.Request) {
	var ctag createTag
	err := json.NewDecoder(r.Body).Decode(&ctag)
	if err != nil {
		a.log.Warnf("warn: create tag, decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	err = a.tags.CreateTag(r.Context(), domain.CreateTag{
		Name: ctag.Name,
		Slug: ctag.Slug,
	})
	if err != nil {
		a.log.Errorf("error: create tag: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func (a *API) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var utag updateTag
	err := json.NewDecoder(r.Body).Decode(&utag)
	if err != nil {
		a.log.Warnf("warn: tag update, decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "tagID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	utag.ID = int(id)
	originalTag, err := a.tags.GetTagByID(r.Context(), utag.ID)
	if err == service.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		a.log.Errorf("error: update tag: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	updTag := domain.UpdateTag{ID: utag.ID}
	if utag.Name != nil {
		updTag.Name = *utag.Name
	} else {
		updTag.Name = originalTag.Name
	}
	if utag.Slug != nil {
		updTag.Slug = *utag.Slug
	} else {
		updTag.Slug = originalTag.Slug
	}
	err = a.tags.UpdateTag(r.Context(), updTag)
	if err != nil {
		a.log.Errorf("error: update tag: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func (a *API) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "tagID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	_, err = a.tags.GetTagByID(r.Context(), int(id))
	if err == service.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		a.log.Errorf("error: delete tag: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	err = a.tags.DeleteTag(r.Context(), domain.DeleteTag{ID: int(id)})
	if err != nil {
		a.log.Errorf("error: delete tag: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}
