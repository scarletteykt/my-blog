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
}

type PostTag struct {
	ID     int `db:"id"`
	TagID  int `db:"tag_id"`
	PostID int `db:"post_id"`
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
