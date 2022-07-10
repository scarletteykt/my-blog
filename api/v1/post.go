package v1

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/scraletteykt/my-blog/internal/post"
	"github.com/scraletteykt/my-blog/internal/tag"
	"github.com/scraletteykt/my-blog/pkg/auth"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"github.com/scraletteykt/my-blog/pkg/server"
	"net/http"
	"strconv"
)

const (
	pageQueryKey = "page"
	postsOnPage  = 30
)

type createPost struct {
	ReadingTime int    `json:"reading_time"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	ImageURL    string `json:"image_url"`
	Content     string `json:"content"`
	Slug        string `json:"slug"`
	TagIDs      []int  `json:"tags"`
}

type updatePost struct {
	ID          int     `json:"id"`
	ReadingTime *int    `json:"reading_time"`
	Publish     *bool   `json:"publish"`
	Title       *string `json:"title"`
	Subtitle    *string `json:"subtitle"`
	ImageURL    *string `json:"image_url"`
	Content     *string `json:"content"`
	Slug        *string `json:"slug"`
	TagIDs      *[]int  `json:"tags"`
}

func (a *API) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	p, err := a.posts.GetPostByID(r.Context(), int(id))
	if err == post.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, p)
}

func (a *API) GetPosts(w http.ResponseWriter, r *http.Request) {
	var page int64
	if p := r.URL.Query().Get(pageQueryKey); p != "" {
		page, _ = strconv.ParseInt(p, 10, 0)
	}
	if page < 0 {
		server.ErrorJSON(w, r, http.StatusBadRequest, errors.New("page param must be positive"))
		return
	}
	if page == 0 {
		page = 1
	}
	posts, err := a.posts.GetPosts(r.Context(), postsOnPage, int((page-1)*postsOnPage))
	if err == post.ErrNotFound {
		server.ResponseJSONWithCode(w, r, http.StatusNoContent, struct{}{})
		return
	}
	if err != nil {
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, posts)
}

func (a *API) CreatePost(w http.ResponseWriter, r *http.Request) {
	var cpost createPost
	err := json.NewDecoder(r.Body).Decode(&cpost)
	if err != nil {
		logger.Warnf("post create: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	u := auth.FromContext(r.Context())
	if u.ID < 0 {
		logger.Warnf("post create: unauthorized access error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusUnauthorized, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&cpost)
	if err != nil {
		logger.Warnf("post create: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	err = a.posts.CreatePost(r.Context(), post.CreatePost{
		UserID:      u.ID,
		ReadingTime: cpost.ReadingTime,
		Title:       cpost.Title,
		Subtitle:    cpost.Subtitle,
		ImageURL:    cpost.ImageURL,
		Content:     cpost.Content,
		Slug:        cpost.Slug,
		TagIDs:      cpost.TagIDs,
	})
	if err != nil {
		logger.Warnf("post create: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func (a *API) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var input updatePost
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	u := auth.FromContext(r.Context())
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Warnf("post update: decoder error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	input.ID = int(id)
	originalPost, err := a.posts.GetPostByID(r.Context(), input.ID)
	if err == post.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if err != nil {
		logger.Warnf("post update: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	if u.ID != originalPost.UserID {
		logger.Warnf("post update: unauthorized access error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusUnauthorized, err)
		return
	}
	updPost := post.UpdatePost{}
	if input.ReadingTime != nil {
		updPost.ReadingTime = *input.ReadingTime
	} else {
		updPost.ReadingTime = originalPost.ReadingTime
	}
	if input.Title != nil {
		updPost.Title = *input.Title
	} else {
		updPost.Title = originalPost.Title
	}
	if input.Subtitle != nil {
		updPost.Subtitle = *input.Subtitle
	} else {
		updPost.Subtitle = originalPost.Subtitle
	}
	if input.ImageURL != nil {
		updPost.ImageURL = *input.ImageURL
	} else {
		updPost.ImageURL = originalPost.ImageURL
	}
	if input.Content != nil {
		updPost.Content = *input.Content
	} else {
		updPost.Content = originalPost.Content
	}
	if input.Slug != nil {
		updPost.Slug = *input.Slug
	} else {
		updPost.Slug = originalPost.Slug
	}
	if input.TagIDs != nil {
		updPost.TagIDs = *input.TagIDs
	} else {
		updPost.TagIDs = getTagIDs(originalPost.Tags)
	}
	if input.Publish != nil {
		if *input.Publish {
			updPost.Status = post.PostStatusPublished
		} else {
			updPost.Status = post.PostStatusDraft
		}
	} else {
		updPost.Status = originalPost.Status
	}

	err = a.posts.UpdatePost(r.Context(), updPost)
	if err != nil {
		logger.Warnf("post update: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func (a *API) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 0)
	if err != nil {
		server.ErrorJSON(w, r, http.StatusBadRequest, err)
		return
	}
	u := auth.FromContext(r.Context())
	originalPost, err := a.posts.GetPostByID(r.Context(), int(id))
	if err == post.ErrNotFound {
		server.ErrorJSON(w, r, http.StatusNotFound, err)
		return
	}
	if u.ID != originalPost.UserID {
		logger.Warnf("post delete: unauthorized access error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusUnauthorized, err)
		return
	}
	err = a.posts.DeletePost(r.Context(), post.DeletePost{ID: int(id)})
	if err != nil {
		logger.Warnf("post delete: error: %s", err.Error())
		server.ErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	server.ResponseJSON(w, r, "ok")
}

func getTagIDs(tags []*tag.Tag) []int {
	out := make([]int, 0)
	for _, t := range tags {
		out = append(out, t.ID)
	}
	return out
}