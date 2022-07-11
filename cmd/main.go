package main

import (
	_ "github.com/lib/pq"
	apiv1 "github.com/scraletteykt/my-blog/api/v1"
	"github.com/scraletteykt/my-blog/internal/config"
	"github.com/scraletteykt/my-blog/internal/post"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/internal/tag"
	"github.com/scraletteykt/my-blog/internal/user"
	"github.com/scraletteykt/my-blog/pkg/server"
	"github.com/scraletteykt/my-blog/pkg/storage"
	"log"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	s, err := storage.New(storage.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})

	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}

	repo := repository.New(s)
	users := user.New(repo)
	posts := post.New(repo, repo, repo)
	tags := tag.New(repo)
	api := apiv1.New(cfg, users, posts, tags)
	srv := server.New()

	if err := srv.Run(cfg, api.Router()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
