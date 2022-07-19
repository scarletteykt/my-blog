package main

import (
	_ "github.com/lib/pq"
	apiv1 "github.com/scraletteykt/my-blog/api/v1"
	"github.com/scraletteykt/my-blog/internal/config"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/internal/service"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"github.com/scraletteykt/my-blog/pkg/server"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.NewConfig(log)

	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	dbConn, err := repository.NewPostgresDB(config.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
		Password: cfg.Postgres.Password,
	})
	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}

	repo := repository.NewRepositories(dbConn, log)
	users := service.NewUsersService(*repo.Users, log)
	posts := service.NewPostsService(*repo.Posts, *repo.Tags, *repo.PostsTags, log)
	tags := service.NewTagsService(*repo.Tags, log)
	api := apiv1.NewAPI(cfg, users, posts, tags, log)
	srv := server.NewServer()

	if err := srv.Run(cfg, api.Router()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
