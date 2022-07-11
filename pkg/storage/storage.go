package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("not found rows in result set")

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Storage struct {
	*sqlx.DB
}

func New(cfg Config) (*Storage, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (s *Storage) Transaction(ctx context.Context, t func(tx *sqlx.Tx) error) error {
	tx, err := s.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	err = t(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.New(fmt.Sprintf("error rollback: %v", txErr))
		}
		return err
	}
	return tx.Commit()
}
