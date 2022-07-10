package repository

import (
	"github.com/scraletteykt/my-blog/pkg/storage"
)

type Repo struct {
	db *storage.Storage
}

func New(db *storage.Storage) *Repo {
	return &Repo{db: db}
}
