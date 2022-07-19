package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

const usersTable = "users"

type UsersRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewUsersRepo(db *sqlx.DB, log logger.Logger) *UsersRepo {
	return &UsersRepo{
		db:  db,
		log: log,
	}
}

func (r *UsersRepo) CreateUser(ctx context.Context, user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", usersTable)

	row := r.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersRepo) GetUser(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE username=$1", usersTable)
	err := r.db.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (r *UsersRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}
