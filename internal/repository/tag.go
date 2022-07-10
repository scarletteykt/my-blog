package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

const (
	tagsTable      = "tags"
	tagSelectQuery = "t.id AS t_id, t.name AS t_name, t.slug AS t_slug"
)

type Tags interface {
	GetTagById(ctx context.Context, id int) (*Tag, error)
	GetTags(ctx context.Context) ([]*Tag, error)
	CreateTag(ctx context.Context, createTag CreateTag) (int, error)
	UpdateTag(ctx context.Context, updateTag UpdateTag) error
	DeleteTag(ctx context.Context, deleteTag DeleteTag) error
}

type Tag struct {
	ID   int    `db:"t_id"`
	Name string `db:"t_name"`
	Slug string `db:"t_slug"`
}

type CreateTag struct {
	Name string
	Slug string
}

type UpdateTag struct {
	ID   int
	Name string
	Slug string
}

type DeleteTag struct {
	ID int
}

func (r *Repo) GetTagById(ctx context.Context, id int) (*Tag, error) {
	query := fmt.Sprintf("SELECT %s FROM %s t WHERE t.id = $1", tagSelectQuery, tagsTable)
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out, err := scanTagRows(rows)
	if err == sql.ErrNoRows || len(out) == 0 {
		return nil, storage.ErrNotFound
	}
	return out[0], err
}

func (r *Repo) GetTags(ctx context.Context) ([]*Tag, error) {
	query := fmt.Sprintf("SELECT %s FROM %s t", tagSelectQuery, tagsTable)
	rows, err := r.db.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out, err := scanTagRows(rows)
	if err == sql.ErrNoRows || len(out) == 0 {
		return nil, storage.ErrNotFound
	}
	return out, nil
}

func (r *Repo) CreateTag(ctx context.Context, createTag CreateTag) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, slug) VALUES($1, $2) RETURNING id", tagsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		rows, err := tx.QueryContext(ctx, query,
			createTag.Name,
			createTag.Slug,
		)
		if err != nil {
			return err
		}
		defer func() { _ = rows.Close() }()
		if rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			if err := rows.Scan(&id); err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *Repo) UpdateTag(ctx context.Context, updateTag UpdateTag) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1, slug=$2 WHERE id=$3", tagsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query,
			updateTag.Name,
			updateTag.Slug,
			updateTag.ID,
		)
		return err
	})
	return err
}

func (r *Repo) DeleteTag(ctx context.Context, deleteTag DeleteTag) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tagsTable)
	_, err := r.db.ExecContext(ctx, query, deleteTag.ID)
	if err != nil {
		return err
	}
	return nil
}

func scanTagRows(rows *sql.Rows) ([]*Tag, error) {
	out := make([]*Tag, 0)
	for rows.Next() {
		t := &Tag{}
		if err := rows.Scan(t); err != nil {
			return nil, err
		}
		if t.ID > 0 {
			out = append(out, t)
		}
	}
	return out, nil
}
