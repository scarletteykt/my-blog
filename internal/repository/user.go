package repository

import (
	"database/sql"
	"fmt"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

const usersTable = "users"

type User struct {
	ID           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

type Users interface {
	CreateUser(createUser CreateUser) (int, error)
	GetUser(username string) (*User, error)
	GetUserByID(id int) (*User, error)
}

type CreateUser struct {
	Username     string
	PasswordHash string
}

func (r *Repo) CreateUser(createUser CreateUser) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", usersTable)

	row := r.db.QueryRow(query, createUser.Username, createUser.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) GetUser(username string) (*User, error) {
	var u User
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE username=$1", usersTable)
	err := r.db.Get(&u, query, username)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &u, err
}

func (r *Repo) GetUserByID(id int) (*User, error) {
	var u User
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&u, query, id)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &u, err
}
