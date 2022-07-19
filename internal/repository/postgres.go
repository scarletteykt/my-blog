package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/scraletteykt/my-blog/internal/config"
)

func NewPostgresDB(cfg config.PostgresConfig) (*sqlx.DB, error) {
	dbConn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
