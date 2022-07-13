package repository

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

const usersTable = "users"

type User struct {
	ID           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

type CreateUser struct {
	Username     string
	PasswordHash string
}

type Users interface {
	CreateUser(ctx context.Context, createUser CreateUser) (int, error)
	GetUser(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, userID int) (*User, error)
}

func (r *Repo) CreateUser(ctx context.Context, createUser CreateUser) (int, error) {
	var id int
	query, args, _ := squirrel.Insert(usersTable).
		SetMap(map[string]interface{}{
			"username":      createUser.Username,
			"password_hash": createUser.PasswordHash,
		}).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	row := r.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) GetUser(ctx context.Context, username string) (*User, error) {
	var u User
	query, args, _ := squirrel.Select(`
			id,
			username,
			password_hash
		`).
		From(usersTable).
		Where("username = ?", username).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	err := r.db.GetContext(ctx, &u, query, args...)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &u, err
}

func (r *Repo) GetUserByID(ctx context.Context, id int) (*User, error) {
	var u User
	query, args, _ := squirrel.Select(`
			id,
			username,
			password_hash
		`).
		From(usersTable).
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	err := r.db.GetContext(ctx, &u, query, args...)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &u, err
}
