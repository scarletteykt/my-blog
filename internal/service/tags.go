package service

import (
	"context"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

type TagsService struct {
	tagsRepo repository.TagsRepo
	log      logger.Logger
}

func NewTagsService(repo repository.TagsRepo, log logger.Logger) *TagsService {
	return &TagsService{
		tagsRepo: repo,
		log:      log,
	}
}

func (t *TagsService) GetTags(ctx context.Context) ([]*domain.Tag, error) {
	dbTags, err := t.tagsRepo.GetTags(ctx)
	out := make([]*domain.Tag, 0)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	for _, dbTag := range dbTags {
		tag := &domain.Tag{
			ID:   dbTag.ID,
			Name: dbTag.Name,
			Slug: dbTag.Slug,
		}
		out = append(out, tag)
	}
	return out, nil
}

func (t *TagsService) GetTagByID(ctx context.Context, tagID int) (*domain.Tag, error) {
	dbTag, err := t.tagsRepo.GetTagByID(ctx, tagID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &domain.Tag{
		ID:   dbTag.ID,
		Name: dbTag.Name,
		Slug: dbTag.Slug,
	}, nil
}

func (t *TagsService) CreateTag(ctx context.Context, createTag domain.CreateTag) error {
	_, err := t.tagsRepo.CreateTag(ctx, repository.CreateTag{
		Name: createTag.Name,
		Slug: createTag.Slug,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *TagsService) UpdateTag(ctx context.Context, updateTag domain.UpdateTag) error {
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

func (t *TagsService) DeleteTag(ctx context.Context, deleteTag domain.DeleteTag) error {
	err := t.tagsRepo.DeleteTag(ctx, repository.DeleteTag{ID: deleteTag.ID})
	if err != nil {
		return err
	}
	return nil
}
