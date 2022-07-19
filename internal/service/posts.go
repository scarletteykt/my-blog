package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/auth"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"time"
)

var ErrNotFound = errors.New("not found rows in result set")

type PostsService struct {
	postsRepo     repository.PostsRepo
	tagsRepo      repository.TagsRepo
	postsTagsRepo repository.PostsTagsRepo
	log           logger.Logger
}

func NewPostsService(postsRepo repository.PostsRepo, tagsRepo repository.TagsRepo, postsTagsRepo repository.PostsTagsRepo, log logger.Logger) *PostsService {
	return &PostsService{
		postsRepo:     postsRepo,
		tagsRepo:      tagsRepo,
		postsTagsRepo: postsTagsRepo,
		log:           log,
	}
}

func (p *PostsService) GetPostByID(ctx context.Context, id int) (*domain.Post, error) {
	u := auth.FromContext(ctx)
	posts, err := p.getPosts(ctx, repository.PostCriteria{
		ID:     id,
		UserID: 0,
		Status: 0,
		TagID:  0,
		Limit:  0,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	post := posts[0]
	if post.UserID != u.ID && post.Status != domain.PostStatusPublished {
		return nil, ErrNotFound
	}
	return posts[0], nil
}

func (p *PostsService) GetPosts(ctx context.Context, limit, offset uint64) ([]*domain.Post, error) {
	return p.getPosts(ctx, repository.PostCriteria{
		ID:     0,
		UserID: 0,
		Status: domain.PostStatusPublished,
		TagID:  0,
		Limit:  limit,
		Offset: offset,
	})
}

func (p *PostsService) GetPostsByTag(ctx context.Context, tagID int, limit, offset uint64) ([]*domain.Post, error) {
	return p.getPosts(ctx, repository.PostCriteria{
		ID:     0,
		UserID: 0,
		Status: domain.PostStatusPublished,
		TagID:  tagID,
		Limit:  limit,
		Offset: offset,
	})
}

func (p *PostsService) GetPostsByUser(ctx context.Context, userID int, limit, offset uint64) ([]*domain.Post, error) {
	return p.getPosts(ctx, repository.PostCriteria{
		ID:     0,
		UserID: userID,
		Status: 0,
		TagID:  0,
		Limit:  limit,
		Offset: offset,
	})
}

func (p *PostsService) CreatePost(ctx context.Context, createPost domain.CreatePost) error {
	postID, err := p.postsRepo.CreatePost(ctx, repository.CreatePost{
		UserID:      createPost.UserID,
		ReadingTime: createPost.ReadingTime,
		Status:      domain.PostStatusDraft,
		Title:       createPost.Title,
		Subtitle:    createPost.Subtitle,
		ImageURL:    createPost.ImageURL,
		Content:     createPost.Content,
		Slug:        createPost.Content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return err
	}
	for _, tagID := range createPost.TagIDs {
		err := p.postsTagsRepo.TagPost(ctx, tagID, postID)
		if err != nil {
			continue
		}
	}
	return nil
}

func (p *PostsService) UpdatePost(ctx context.Context, updatePost domain.UpdatePost) error {
	var publishedAt sql.NullTime
	if updatePost.Status == domain.PostStatusPublished {
		publishedAt.Time = time.Now()
		publishedAt.Valid = true
	} else {
		publishedAt.Time = time.Time{}
		publishedAt.Valid = false
	}
	err := p.postsRepo.UpdatePost(ctx, repository.UpdatePost{
		ID:          updatePost.ID,
		ReadingTime: updatePost.ReadingTime,
		Status:      updatePost.Status,
		Title:       updatePost.Title,
		Subtitle:    updatePost.Subtitle,
		ImageURL:    updatePost.ImageURL,
		Content:     updatePost.Content,
		Slug:        updatePost.Slug,
		PublishedAt: publishedAt,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return err
	}
	err = p.postsTagsRepo.UpdatePostTags(ctx, updatePost.TagIDs, updatePost.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostsService) DeletePost(ctx context.Context, deletePost domain.DeletePost) error {
	err := p.postsRepo.DeletePost(ctx, repository.DeletePost{
		ID:        deletePost.ID,
		Status:    domain.PostStatusDeleted,
		DeletedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostsService) getPosts(ctx context.Context, criteria repository.PostCriteria) ([]*domain.Post, error) {
	dbPosts, err := p.postsRepo.GetPostsByCriteria(ctx, criteria)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	if dbPosts == nil || len(dbPosts) == 0 {
		return nil, ErrNotFound
	}

	posts := make([]*domain.Post, 0)

	for _, dbPost := range dbPosts {
		var (
			publishedAt time.Time
			deletedAt   time.Time
		)
		if dbPost.PublishedAt.Valid {
			publishedAt = dbPost.PublishedAt.Time
		}
		if dbPost.DeletedAt.Valid {
			deletedAt = dbPost.DeletedAt.Time
		}
		tags := make([]*domain.Tag, 0)
		pst := &domain.Post{
			ID:          dbPost.ID,
			UserID:      dbPost.UserID,
			ReadingTime: dbPost.ReadingTime,
			Status:      dbPost.Status,
			Title:       dbPost.Title,
			Subtitle:    dbPost.Subtitle,
			ImageURL:    dbPost.ImageURL,
			Content:     dbPost.Content,
			Slug:        dbPost.Slug,
			PublishedAt: publishedAt,
			CreatedAt:   dbPost.CreatedAt,
			UpdatedAt:   dbPost.UpdatedAt,
			DeletedAt:   deletedAt,
			Tags:        tags,
		}
		for _, dbTag := range dbPost.Tags {
			t := &domain.Tag{
				ID:   dbTag.ID,
				Name: dbTag.Name,
				Slug: dbTag.Slug,
			}
			pst.Tags = append(pst.Tags, t)
		}
		posts = append(posts, pst)
	}
	return posts, nil
}
