package user

import (
	"context"
	"errors"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/storage"
)

var (
	ErrNotFound          = errors.New("wrong username or password")
	ErrUserAlreadyExists = errors.New("user with given username already exists")
)

type Users struct {
	repo *repository.Repo
}

func New(r *repository.Repo) *Users {
	return &Users{
		repo: r,
	}
}

func (s *Users) GetUser(ctx context.Context, username string) (*User, error) {
	userDB, err := s.repo.GetUser(ctx, username)
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:           userDB.ID,
		Username:     userDB.Username,
		PasswordHash: userDB.PasswordHash,
	}

	return u, nil
}

func (s *Users) CreateUser(ctx context.Context, createUser CreateUser) (*User, error) {
	_, err := s.repo.GetUser(ctx, createUser.Username)
	if err == storage.ErrNotFound {
		id, err := s.repo.CreateUser(ctx, repository.CreateUser{
			Username:     createUser.Username,
			PasswordHash: createUser.PasswordHash,
		})
		if err != nil {
			return nil, err
		}

		userDB, err := s.repo.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
		u := &User{
			ID:           userDB.ID,
			Username:     userDB.Username,
			PasswordHash: userDB.PasswordHash,
		}

		return u, nil
	} else if err != nil {
		return nil, err
	} else {
		return nil, ErrUserAlreadyExists
	}
}
