package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

const (
	postsTable      = "posts"
	postSelectQuery = "p.id AS p_id, p.user_id AS p_user_id, p.reading_time AS p_reading_time, p.status AS p_status, " +
		"p.title AS p_title, p.subtitle AS p_subtitle, p.image_url AS p_image_url, p.content AS p_content, p.slug AS p_slug, " +
		"p.published_at AS p_published_at, p.created_at AS p_created_at, p.updated_at AS p_updated_at, p.deleted_at AS p_deleted_at, " +
		"t.id AS t_id, t.name AS t_name, t.slug AS t_slug "
	postsPerPage = 30
)

type Posts interface {
	GetPostsByCriteria(ctx context.Context, criteria PostCriteria) ([]*Post, error)
	CreatePost(ctx context.Context, createPost CreatePost) (int, error)
	UpdatePost(ctx context.Context, updatePost UpdatePost) error
	DeletePost(ctx context.Context, deletePost DeletePost) error
}

type Post struct {
	ID          int       `db:"p_id"`
	UserID      int       `db:"p_user_id"`
	ReadingTime int       `db:"p_reading_time"`
	Status      int       `db:"p_status"`
	Title       string    `db:"p_title"`
	Subtitle    string    `db:"p_subtitle"`
	ImageURL    string    `db:"p_image_url"`
	Content     string    `db:"p_content"`
	Slug        string    `db:"p_slug"`
	PublishedAt NullTime  `db:"p_published_at"`
	CreatedAt   time.Time `db:"p_created_at"`
	UpdatedAt   time.Time `db:"p_updated_at"`
	DeletedAt   NullTime  `db:"p_deleted_at"`
	Tags        []*Tag
}

type PostCriteria struct {
	ID     int
	UserID int
	Status int
	TagID  int
	Limit  int
	Offset int
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
	PublishedAt NullTime
	UpdatedAt   time.Time
}

type DeletePost struct {
	ID        int
	Status    int
	DeletedAt time.Time
}

func (r *Repo) GetPostsByCriteria(ctx context.Context, criteria PostCriteria) ([]*Post, error) {
	where := "1=1 "
	if criteria.UserID > 0 {
		where += "AND p.user_id = :user_id "
	}
	if criteria.Status != PostStatusDeleted {
		where += fmt.Sprintf("AND p.status <> %d ", PostStatusDeleted)
	}
	if criteria.Status > 0 {
		where += "AND p.status = :status "
	}
	if criteria.ID > 0 {
		where += "AND p.id = :id "
	}
	if criteria.TagID > 0 {
		where += "AND t.id = :tag_id "
	}
	if criteria.Limit == 0 {
		criteria.Limit = postsPerPage
	}

	subquery := fmt.Sprintf("SELECT p.id FROM %s p LEFT JOIN %s pt ON p.id = pt.post_id LEFT JOIN %s t ON pt.tag_id = t.id  "+
		"WHERE %s GROUP BY p.id LIMIT %d OFFSET %d",
		postsTable, postsTagsTable, tagsTable, where, criteria.Limit, criteria.Offset)
	query := fmt.Sprintf("SELECT %s FROM %s p LEFT JOIN %s pt ON p.id = pt.post_id LEFT JOIN %s t ON pt.tag_id = t.id WHERE p.id IN (%s) "+
		"ORDER BY p.created_at DESC ",
		postSelectQuery, postsTable, postsTagsTable, tagsTable, subquery)

	rows, err := r.db.NamedQueryContext(ctx, query, criteria)
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
	query := fmt.Sprintf("INSERT INTO %s "+
		"(user_id, reading_time, status, title, subtitle, image_url, content, slug, created_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", postsTable)

	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		rows, err := tx.QueryContext(ctx, query,
			createPost.UserID,
			createPost.ReadingTime,
			createPost.Status,
			createPost.Title,
			createPost.Subtitle,
			createPost.ImageURL,
			createPost.Content,
			createPost.Slug,
			createPost.CreatedAt,
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

func (r *Repo) UpdatePost(ctx context.Context, updatePost UpdatePost) error {
	query := fmt.Sprintf("UPDATE %s SET "+
		"reading_time=$1, status=$2, title=$3, subtitle=$4, image_url=$5, content=$6, slug=$7, published_at=$8, updated_at=$9 "+
		"WHERE id=$10", postsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query,
			updatePost.ReadingTime,
			updatePost.Status,
			updatePost.Title,
			updatePost.Subtitle,
			updatePost.ImageURL,
			updatePost.Content,
			updatePost.Slug,
			updatePost.PublishedAt,
			updatePost.UpdatedAt,
			updatePost.ID,
		)
		return err
	})
	return err
}

func (r *Repo) DeletePost(ctx context.Context, deletePost DeletePost) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1, deleted_at=$2 WHERE id=$3", postsTable)
	err := r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query,
			deletePost.Status,
			deletePost.DeletedAt,
			deletePost.ID,
		)
		return err
	})
	return err
}

func scanPostRows(rows *sqlx.Rows) ([]*Post, error) {
	posts := make(map[int]*Post)
	out := make([]*Post, 0)
	for rows.Next() {
		p := &Post{}
		t := &Tag{}
		if err := rows.Scan(p, t); err != nil {
			return nil, err
		}
		if _, ok := posts[p.ID]; !ok {
			posts[p.ID] = p
		}
		if t.ID > 0 {
			posts[p.ID].Tags = append(posts[p.ID].Tags, t)
		}
	}
	for _, p := range posts {
		out = append(out, p)
	}
	return out, nil
}
