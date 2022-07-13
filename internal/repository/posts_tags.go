package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const postsTagsTable = "posts_tags"

type PostsTags interface {
	TagPost(ctx context.Context, tagID, postID int) error
	UntagPost(ctx context.Context, tagID, postID int) error
	UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error
}

func (r *Repo) TagPost(ctx context.Context, tagID, postID int) error {
	query, args, _ := squirrel.Insert(postsTagsTable).
		SetMap(map[string]interface{}{
			"tag_id":  tagID,
			"post_id": postID,
		}).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) UntagPost(ctx context.Context, tagID, postID int) error {
	query, args, _ := squirrel.Delete(postsTagsTable).
		Where("tag_id = ? AND post_id = ?", tagID, postID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error {
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		query, args, _ := squirrel.Delete(postsTagsTable).
			Where("post_id = ?", postID).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		query, _, _ = squirrel.Insert(postsTagsTable).
			PlaceholderFormat(squirrel.Dollar).
			Columns("tag_id", "post_id").
			Values("", "").
			ToSql()
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			_, err = stmt.ExecContext(ctx, tagID, postID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
