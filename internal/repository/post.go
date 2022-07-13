package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

const (
	postsTable   = "posts"
	postsPerPage = 30
)

type Posts interface {
	GetPostsByCriteria(ctx context.Context, criteria PostCriteria) ([]*Post, error)
	CreatePost(ctx context.Context, createPost CreatePost) (int, error)
	UpdatePost(ctx context.Context, updatePost UpdatePost) error
	DeletePost(ctx context.Context, deletePost DeletePost) error
}

type Post struct {
	ID          int          `db:"p_id"`
	UserID      int          `db:"p_user_id"`
	ReadingTime int          `db:"p_reading_time"`
	Status      int          `db:"p_status"`
	Title       string       `db:"p_title"`
	Subtitle    string       `db:"p_subtitle"`
	ImageURL    string       `db:"p_image_url"`
	Content     string       `db:"p_content"`
	Slug        string       `db:"p_slug"`
	PublishedAt sql.NullTime `db:"p_published_at"`
	CreatedAt   time.Time    `db:"p_created_at"`
	UpdatedAt   time.Time    `db:"p_updated_at"`
	DeletedAt   sql.NullTime `db:"p_deleted_at"`
	Tags        []*Tag
}

type PostTag struct {
	ID          int            `db:"p_id"`
	UserID      int            `db:"p_user_id"`
	ReadingTime int            `db:"p_reading_time"`
	Status      int            `db:"p_status"`
	Title       string         `db:"p_title"`
	Subtitle    string         `db:"p_subtitle"`
	ImageURL    string         `db:"p_image_url"`
	Content     string         `db:"p_content"`
	Slug        string         `db:"p_slug"`
	PublishedAt sql.NullTime   `db:"p_published_at"`
	CreatedAt   time.Time      `db:"p_created_at"`
	UpdatedAt   time.Time      `db:"p_updated_at"`
	DeletedAt   sql.NullTime   `db:"p_deleted_at"`
	TagID       sql.NullInt32  `db:"t_id"`
	TagName     sql.NullString `db:"t_name"`
	TagSlug     sql.NullString `db:"t_slug"`
}

type PostCriteria struct {
	ID     int `db:"p_id"`
	UserID int `db:"p_user_id"`
	Status int `db:"p_status"`
	TagID  int `db:"t_id"`
	Limit  uint64
	Offset uint64
}

type CreatePost struct {
	UserID      int
	ReadingTime int
	Status      int
	Title       string
	Subtitle    string
	ImageURL    string
	Content     string
	Slug        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdatePost struct {
	ID          int
	ReadingTime int
	Status      int
	Title       string
	Subtitle    string
	ImageURL    string
	Content     string
	Slug        string
	PublishedAt sql.NullTime
	UpdatedAt   time.Time
}

type DeletePost struct {
	ID        int
	Status    int
	DeletedAt time.Time
}

func (r *Repo) GetPostsByCriteria(ctx context.Context, criteria PostCriteria) ([]*Post, error) {
	sb := squirrel.Select(`p.id`).
		From(postsTable + " p").
		LeftJoin(postsTagsTable + " pt ON p.id = pt.post_id").
		LeftJoin(tagsTable + " t ON pt.tag_id = t.id")

	if criteria.UserID > 0 {
		sb = sb.Where("p.user_id = :p_user_id", criteria.UserID)
	}
	if criteria.Status != PostStatusDeleted {
		sb = sb.Where(fmt.Sprintf("p.status <> %d", PostStatusDeleted))
	}
	if criteria.Status > 0 {
		sb = sb.Where("p.status = :p_status", criteria.Status)
	}
	if criteria.ID > 0 {
		sb = sb.Where("p.id = :p_id", criteria.ID)
	}
	if criteria.TagID > 0 {
		sb = sb.Where("t.id = :t_id", criteria.TagID)
	}
	if criteria.Limit == 0 {
		criteria.Limit = postsPerPage
	}

	sb.GroupBy("p.id").Limit(criteria.Limit).Offset(criteria.Offset)

	query, _, _ := squirrel.Select(`
			p.id AS p_id, 
			p.user_id AS p_user_id, 
			p.reading_time AS p_reading_time, 
			p.status AS p_status,
			p.title AS p_title, 
			p.subtitle AS p_subtitle, 
			p.image_url AS p_image_url, 
			p.content AS p_content, 
			p.slug AS p_slug, 
			p.published_at AS p_published_at, 
			p.created_at AS p_created_at, 
			p.updated_at AS p_updated_at, 
			p.deleted_at AS p_deleted_at,
			t.id AS t_id, 
			t.name AS t_name, 
			t.slug AS t_slug`).
		From(postsTable + " p").
		LeftJoin(postsTagsTable + " pt ON p.id = pt.post_id").
		LeftJoin(tagsTable + " t ON pt.tag_id = t.id").
		Where(subquery("p.id IN", sb)).
		OrderBy("p.created_at DESC").
		ToSql()

	rows, err := r.db.NamedQueryContext(ctx, query, criteria)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out, err := scanPostRows(rows)
	if err == sql.ErrNoRows || len(out) == 0 {
		return nil, storage.ErrNotFound
	}
	return out, nil
}

func (r *Repo) CreatePost(ctx context.Context, createPost CreatePost) (int, error) {
	var id int
	query, args, _ := squirrel.Insert(postsTable).
		SetMap(map[string]interface{}{
			"user_id":      createPost.UserID,
			"reading_time": createPost.ReadingTime,
			"status":       createPost.Status,
			"title":        createPost.Title,
			"subtitle":     createPost.Subtitle,
			"image_url":    createPost.ImageURL,
			"content":      createPost.Content,
			"slug":         createPost.Slug,
			"created_at":   createPost.CreatedAt,
			"updated_at":   createPost.UpdatedAt,
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

func (r *Repo) UpdatePost(ctx context.Context, updatePost UpdatePost) error {
	query, args, _ := squirrel.Update(postsTable).
		SetMap(map[string]interface{}{
			"reading_time": updatePost.ReadingTime,
			"status":       updatePost.Status,
			"title":        updatePost.Title,
			"subtitle":     updatePost.Subtitle,
			"image_url":    updatePost.ImageURL,
			"content":      updatePost.Content,
			"slug":         updatePost.Slug,
			"published_at": updatePost.PublishedAt,
			"updated_at":   updatePost.UpdatedAt,
		}).
		Where("id = ?", updatePost.ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) DeletePost(ctx context.Context, deletePost DeletePost) error {
	query, args, _ := squirrel.Update(postsTable).
		SetMap(map[string]interface{}{
			"status":     deletePost.Status,
			"deleted_at": deletePost.DeletedAt,
		}).
		Where("id = ?", deletePost.ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func scanPostRows(rows *sqlx.Rows) ([]*Post, error) {
	posts := make(map[int]*Post)
	out := make([]*Post, 0)
	for rows.Next() {
		pt := &PostTag{}
		if err := rows.StructScan(pt); err != nil {
			return nil, err
		}
		if _, ok := posts[pt.ID]; !ok {
			p := &Post{
				ID:          pt.ID,
				UserID:      pt.UserID,
				ReadingTime: pt.ReadingTime,
				Status:      pt.Status,
				Title:       pt.Title,
				Subtitle:    pt.Subtitle,
				ImageURL:    pt.ImageURL,
				Content:     pt.Content,
				Slug:        pt.Slug,
				PublishedAt: pt.PublishedAt,
				CreatedAt:   pt.CreatedAt,
				UpdatedAt:   pt.UpdatedAt,
				DeletedAt:   pt.DeletedAt,
			}
			p.Tags = make([]*Tag, 0)
			posts[p.ID] = p
		}
		if pt.TagID.Valid && pt.TagName.Valid && pt.TagSlug.Valid {
			t := &Tag{
				ID:   int(pt.TagID.Int32),
				Name: pt.TagName.String,
				Slug: pt.TagSlug.String,
			}
			posts[pt.ID].Tags = append(posts[pt.ID].Tags, t)
		}
	}
	for _, p := range posts {
		out = append(out, p)
	}
	return out, nil
}

func subquery(prefix string, sb squirrel.SelectBuilder) squirrel.Sqlizer {
	s, params, _ := sb.ToSql()
	return squirrel.Expr(prefix+" ("+s+")", params)
}
