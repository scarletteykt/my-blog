package main

import (
	"github.com/scraletteykt/my-blog/internal/post"
	"github.com/scraletteykt/my-blog/internal/tag"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	apiv1 "github.com/scraletteykt/my-blog/api/v1"
	"github.com/scraletteykt/my-blog/internal/config"
	"github.com/scraletteykt/my-blog/internal/repository"
	"github.com/scraletteykt/my-blog/internal/user"
	"github.com/scraletteykt/my-blog/pkg/server"
	"github.com/scraletteykt/my-blog/pkg/storage"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
	}

	s, err := storage.New(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}

	repo := repository.New(s)
	users := user.New(repo)
	posts := post.New(repo, repo, repo)
	tags := tag.New(repo)
	api := apiv1.New(users, posts, tags)
	srv := server.New()

	if err := srv.Run(nil, api.Router()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
