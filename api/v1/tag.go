package v1

import "net/http"

type createTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type updateTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type deleteTag struct {
	ID int `json:"id"`
}

func (a *API) CreateTag(w http.ResponseWriter, r *http.Request) {

}

func (a *API) UpdateTag(w http.ResponseWriter, r *http.Request) {
	
}

func (a *API) DeleteTag(w http.ResponseWriter, r *http.Request) {

}
