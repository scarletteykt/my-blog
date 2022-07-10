package post

import (
	"context"
	"github.com/pkg/errors"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/storage"
	"time"
)

var ErrNotFound = errors.New("not found rows in result set")

type Posts struct {
	postsRepo     repository.Posts
	tagsRepo      repository.Tags
	postsTagsRepo repository.PostsTags
}

func New(postsRepo repository.Posts, tagsRepo repository.Tags, postsTagsRepo repository.PostsTags) *Posts {
	return &Posts{
		postsRepo:     postsRepo,
		tagsRepo:      tagsRepo,
		postsTagsRepo: postsTagsRepo,
	}
}

func (p *Posts) GetPostByID(ctx context.Context, id int) (*Post, error) {
	posts, err := p.getPosts(ctx, repository.PostCriteria{
		ID:     id,
		UserID: 0,
		Status: PostStatusPublished,
		TagID:  0,
		Limit:  0,
		Offset: 0,
	})
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return posts[0], nil
}

func (p *Posts) GetPosts(ctx context.Context, limit, offset int) ([]*Post, error) {
	return p.getPosts(ctx, repository.PostCriteria{
		ID:     0,
		UserID: 0,
		Status: PostStatusPublished,
		TagID:  0,
		Limit:  limit,
		Offset: offset,
	})
}

func (p *Posts) GetPostsByTag(ctx context.Context, tagID, limit, offset int) ([]*Post, error) {
	return p.getPosts(ctx, repository.PostCriteria{
		ID:     0,
		UserID: 0,
		Status: PostStatusPublished,
		TagID:  tagID,
		Limit:  limit,
		Offset: offset,
	})
}

func (p *Posts) CreatePost(ctx context.Context, createPost CreatePost) error {
	postID, err := p.postsRepo.CreatePost(ctx, repository.CreatePost{
		UserID:      createPost.UserID,
		ReadingTime: createPost.ReadingTime,
		Status:      PostStatusDraft,
		Title:       createPost.Title,
		Subtitle:    createPost.Subtitle,
		ImageURL:    createPost.ImageURL,
		Content:     createPost.Content,
		Slug:        createPost.Content,
		CreatedAt:   time.Now(),
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

func (p *Posts) UpdatePost(ctx context.Context, updatePost UpdatePost) error {
	status := PostStatusDraft
	var publishedAt repository.NullTime
	if updatePost.Publish {
		status = PostStatusPublished
		publishedAt.Time = time.Now()
		publishedAt.Valid = true
	}
	err := p.postsRepo.UpdatePost(ctx, repository.UpdatePost{
		ID:          updatePost.ID,
		ReadingTime: updatePost.ReadingTime,
		Status:      status,
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
	for _, tagID := range updatePost.TagIDs {
		
	}
}

func (p *Posts) getPosts(ctx context.Context, criteria repository.PostCriteria) ([]*Post, error) {
	dbPosts, err := p.postsRepo.GetPostsByCriteria(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if dbPosts == nil || len(dbPosts) == 0 {
		return nil, ErrNotFound
	}

	posts := make([]*Post, 0)

	for _, dbPost := range dbPosts {
		tags := make([]*Tag, 0)
		pst := &Post{
			ID:          dbPost.ID,
			UserID:      dbPost.UserID,
			ReadingTime: dbPost.ReadingTime,
			Status:      dbPost.Status,
			Title:       dbPost.Title,
			Subtitle:    dbPost.Subtitle,
			ImageURL:    dbPost.ImageURL,
			Content:     dbPost.Content,
			Slug:        dbPost.Slug,
			PublishedAt: dbPost.PublishedAt,
			CreatedAt:   dbPost.CreatedAt,
			UpdatedAt:   dbPost.UpdatedAt,
			DeletedAt:   dbPost.DeletedAt,
			Tags:        tags,
		}
		for _, dbTag := range dbPost.Tags {
			tag := &Tag{
				ID:   dbTag.ID,
				Name: dbTag.Name,
				Slug: dbTag.Slug,
			}
			pst.Tags = append(pst.Tags, tag)
		}
		posts = append(posts, pst)
	}
	return posts, nil
}
