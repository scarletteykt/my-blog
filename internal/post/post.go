package post

import (
	"github.com/scraletteykt/my-blog/internal/tag"
	"time"
)

type Post struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	ReadingTime int        `json:"reading_time"`
	Status      int        `json:"status"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle"`
	ImageURL    string     `json:"image_url"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	PublishedAt time.Time  `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   time.Time  `json:"deleted_at"`
	Tags        []*tag.Tag `json:"tags"`
}

type CreatePost struct {
	UserID      int
	ReadingTime int
	Title       string
	Subtitle    string
	ImageURL    string
	Content     string
	Slug        string
	TagIDs      []int
}

type UpdatePost struct {
	ID          int
	ReadingTime int
	Status      int
	Title       string
	Subtitle    string
	ImageURL    string
	Content     string
	Slug        string
	TagIDs      []int
}

type DeletePost struct {
	ID int
}
