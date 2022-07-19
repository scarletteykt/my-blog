package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const postsTagsTable = "posts_tags"

type PostsTags interface {
	TagPost(ctx context.Context, tagID, postID int) error
	UntagPost(ctx context.Context, tagID, postID int) error
	UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error
}

func (r *Repo) TagPost(ctx context.Context, tagID, postID int) error {
	query := fmt.Sprintf("INSERT INTO %s (tag_id, post_id) VALUES ($1, $2)", postsTagsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query, tagID, postID)
		return err
	})
	return err
}

func (r *Repo) UntagPost(ctx context.Context, tagID, postID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE tag_id = $1 AND post_id = $2", postsTagsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query, tagID, postID)
		return err
	})
	return err
}

func (r *Repo) UpdatePostTags(ctx context.Context, tagIDs []int, postID int) error {
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		query := fmt.Sprintf("DELETE FROM %s WHERE post_id = $1", postsTagsTable)
		_, err := tx.ExecContext(ctx, query, postID)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("INSERT INTO %s (tag_id, post_id) VALUES($1, $2)", postsTagsTable)
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
