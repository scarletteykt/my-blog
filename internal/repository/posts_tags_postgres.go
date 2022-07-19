package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"log"
)

const postsTagsTable = "posts_tags"

type PostsTagsRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewPostsTagsRepo(db *sqlx.DB, log logger.Logger) *PostsTagsRepo {
	return &PostsTagsRepo{
		db:  db,
		log: log,
	}
}

func (r *PostsTagsRepo) TagPost(ctx context.Context, tagID, postID int) error {
	query, args, _ := squirrel.Insert(postsTagsTable).
		SetMap(map[string]interface{}{
			"tag_id":  tagID,
			"post_id": postID,
		}).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostsTagsRepo) UntagPost(ctx context.Context, tagID, postID int) error {
	query, args, _ := squirrel.Delete(postsTagsTable).
		Where("tag_id = ? AND post_id = ?", tagID, postID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostsTagsRepo) UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	query, args, _ := squirrel.Delete(postsTagsTable).
		Where("post_id = ?", postID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, execErr := tx.ExecContext(ctx, query, args...)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update tags failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		log.Fatalf("update tags failed: %v", execErr)
		return err
	}
	query, _, _ = squirrel.Insert(postsTagsTable).
		Columns("tag_id", "post_id").
		Values("", "").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	stmt, err := tx.Prepare(query)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update tags failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		log.Fatalf("update tags failed: %v", execErr)
		return err
	}
	defer stmt.Close()
	for _, tagID := range tagIDs {
		_, err = stmt.ExecContext(ctx, tagID, postID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update tags failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
			}
			log.Fatalf("update tags failed: %v", execErr)
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	return nil
}
