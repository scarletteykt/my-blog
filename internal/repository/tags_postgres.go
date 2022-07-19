package repository

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

const tagsTable = "tags"

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

type TagsRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewTagsRepo(db *sqlx.DB, log logger.Logger) *TagsRepo {
	return &TagsRepo{
		db:  db,
		log: log,
	}
}

func (r *TagsRepo) GetTagByID(ctx context.Context, id int) (*Tag, error) {
	query, args, _ := squirrel.Select(`
			t.id AS t_id,
			t.name AS t_name,
			t.slug AS t_slug
		`).
		From(tagsTable+" t").
		Where("t.id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out, err := scanTagRows(rows)
	if err == sql.ErrNoRows || len(out) == 0 {
		return nil, ErrNotFound
	}
	return out[0], err
}

func (r *TagsRepo) GetTags(ctx context.Context) ([]*Tag, error) {
	query, _, _ := squirrel.Select(`
			t.id AS t_id,
			t.name AS t_name,
			t.slug AS t_slug
		`).
		From(tagsTable + " t").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	rows, err := r.db.QueryxContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out, err := scanTagRows(rows)
	if err == sql.ErrNoRows || len(out) == 0 {
		return nil, ErrNotFound
	}
	return out, nil
}

func (r *TagsRepo) CreateTag(ctx context.Context, createTag CreateTag) (int, error) {
	var id int
	query, args, _ := squirrel.Insert(tagsTable).
		SetMap(map[string]interface{}{
			"name": createTag.Name,
			"slug": createTag.Slug,
		}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()
	if rows.Next() {
		if err = rows.Err(); err != nil {
			return 0, err
		}
		if err = rows.Scan(&id); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *TagsRepo) UpdateTag(ctx context.Context, updateTag UpdateTag) error {
	query, args, _ := squirrel.Update(tagsTable).
		SetMap(map[string]interface{}{
			"name": updateTag.Name,
			"slug": updateTag.Slug,
		}).
		Where("id = ?", updateTag.ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TagsRepo) DeleteTag(ctx context.Context, deleteTag DeleteTag) error {
	query, args, _ := squirrel.Delete(tagsTable).
		Where("id = ?", deleteTag.ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func scanTagRows(rows *sqlx.Rows) ([]*Tag, error) {
	out := make([]*Tag, 0)
	for rows.Next() {
		t := &Tag{}
		if err := rows.StructScan(t); err != nil {
			return nil, err
		}
		if t.ID > 0 {
			out = append(out, t)
		}
	}
	return out, nil
}
