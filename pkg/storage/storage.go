package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("not found rows in result set")

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `sslmode:"host"`
}

type Storage struct {
	*sqlx.DB
}

func New(cfg Config) (*Storage, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s, port=%s, user=%s, password=%s,dbname=%s, sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	s := &Storage{}
	s.DB = db

	return s, nil
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
	}
	return nil
}
