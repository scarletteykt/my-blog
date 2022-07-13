package v1

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/scraletteykt/my-blog/internal/tag"
	"github.com/scraletteykt/my-blog/pkg/logger"
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
	if err == tag.ErrNotFound {
		server.ResponseJSONWithCode(w, r, http.StatusNoContent, struct{}{})
		return
	}
	if err != nil {
		logger.Warnf("tag get by id: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, t)
}

func (a *API) GetTags(w http.ResponseWriter, r *http.Request) {
	t, err := a.tags.GetTags(r.Context())
	if err == tag.ErrNotFound {
		server.ResponseJSONWithCode(w, r, http.StatusNoContent, struct{}{})
		return
	}
	if err != nil {
		logger.Warnf("tags get: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, t)
}

func (a *API) CreateTag(w http.ResponseWriter, r *http.Request) {
	var ctag createTag
	err := json.NewDecoder(r.Body).Decode(&ctag)
	if err != nil {
		logger.Warnf("tag create: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	err = a.tags.CreateTag(r.Context(), tag.CreateTag{
		Name: ctag.Name,
		Slug: ctag.Slug,
	})
	if err != nil {
		logger.Warnf("tag create: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func (a *API) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var utag updateTag
	err := json.NewDecoder(r.Body).Decode(&utag)
	if err != nil {
		logger.Warnf("tag update: decoder error: %s", err.Error())
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
	if err == tag.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		logger.Warnf("tag update: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	updTag := tag.UpdateTag{ID: utag.ID}
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
		logger.Warnf("tag update: error: %s", err.Error())
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
	if err == tag.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		logger.Warnf("tag delete: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	err = a.tags.DeleteTag(r.Context(), tag.DeleteTag{ID: int(id)})
	if err != nil {
		logger.Warnf("tag delete: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}
