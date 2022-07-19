package tag

import (
	"context"
	"github.com/pkg/errors"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

var ErrNotFound = errors.New("not found rows in result set")

type Tags struct {
	tagsRepo repository.Tags
}

func New(repo repository.Tags) *Tags {
	return &Tags{
		tagsRepo: repo,
	}
}

func (t *Tags) GetTags(ctx context.Context) ([]*Tag, error) {
	dbTags, err := t.tagsRepo.GetTags(ctx)
	out := make([]*Tag, 0)
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	for _, dbTag := range dbTags {
		tag := &Tag{
			ID:   dbTag.ID,
			Name: dbTag.Name,
			Slug: dbTag.Slug,
		}
		out = append(out, tag)
	}
	return out, nil
}

func (t *Tags) GetTagByID(ctx context.Context, tagID int) (*Tag, error) {
	dbTag, err := t.tagsRepo.GetTagById(ctx, tagID)
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &Tag{
		ID:   dbTag.ID,
		Name: dbTag.Name,
		Slug: dbTag.Slug,
	}, nil
}

func (t *Tags) CreateTag(ctx context.Context, createTag CreateTag) error {
	_, err := t.tagsRepo.CreateTag(ctx, repository.CreateTag{
		Name: createTag.Name,
		Slug: createTag.Slug,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *Tags) UpdateTag(ctx context.Context, updateTag UpdateTag) error {
	err := t.tagsRepo.UpdateTag(ctx, repository.UpdateTag{
		ID:   updateTag.ID,
		Name: updateTag.Name,
		Slug: updateTag.Slug,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *Tags) DeleteTag(ctx context.Context, deleteTag DeleteTag) error {
	err := t.tagsRepo.DeleteTag(ctx, repository.DeleteTag{ID: deleteTag.ID})
	if err != nil {
		return err
	}
	return nil
}
