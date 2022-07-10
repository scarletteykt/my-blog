package auth

import (
	"github.com/scraletteykt/my-blog/internal/service"
)

type Options struct {
	Secret   string
	Services *service.Services
}
