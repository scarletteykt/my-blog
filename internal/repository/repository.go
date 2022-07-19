package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

var ErrNotFound = errors.New("not found rows in result set")

type Users interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
}

type Posts interface {
	GetPostsByCriteria(ctx context.Context, criteria PostCriteria) ([]*Post, error)
	CreatePost(ctx context.Context, createPost CreatePost) (int, error)
	UpdatePost(ctx context.Context, updatePost UpdatePost) error
	DeletePost(ctx context.Context, deletePost DeletePost) error
}

type Tags interface {
	GetTagByID(ctx context.Context, id int) (*Tag, error)
	GetTags(ctx context.Context) ([]*Tag, error)
	CreateTag(ctx context.Context, createTag CreateTag) (int, error)
	UpdateTag(ctx context.Context, updateTag UpdateTag) error
	DeleteTag(ctx context.Context, deleteTag DeleteTag) error
}

type PostsTags interface {
	TagPost(ctx context.Context, tagID, postID int) error
	UntagPost(ctx context.Context, tagID, postID int) error
	UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error
}

type Repositories struct {
	Users     *UsersRepo
	Posts     *PostsRepo
	Tags      *TagsRepo
	PostsTags *PostsTagsRepo
}

func NewRepositories(db *sqlx.DB, log logger.Logger) *Repositories {
	return &Repositories{
		Users:     NewUsersRepo(db, log),
		Posts:     NewPostsRepo(db, log),
		Tags:      NewTagsRepo(db, log),
		PostsTags: NewPostsTagsRepo(db, log),
	}
}
