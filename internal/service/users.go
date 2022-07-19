package service

import (
	"context"
	"errors"
	"github.com/scraletteykt/my-blog/internal/domain"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

var (
	ErrForbidden         = errors.New("wrong username or password")
	ErrUserAlreadyExists = errors.New("user with given username already exists")
)

type UsersService struct {
	repo repository.UsersRepo
	log  logger.Logger
}

func NewUsersService(r repository.UsersRepo, log logger.Logger) *UsersService {
	return &UsersService{
		repo: r,
		log:  log,
	}
}

func (s *UsersService) GetUser(ctx context.Context, username string) (*domain.User, error) {
	userDB, err := s.repo.GetUser(ctx, username)
	if err == repository.ErrNotFound {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           userDB.ID,
		Username:     userDB.Username,
		PasswordHash: userDB.PasswordHash,
	}

	return user, nil
}

func (s *UsersService) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	_, err := s.repo.GetUser(ctx, user.Username)
	if err == repository.ErrNotFound {
		id, err := s.repo.CreateUser(ctx, domain.User{
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
		})
		if err != nil {
			return nil, err
		}

		userDB, err := s.repo.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
		u := &domain.User{
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
