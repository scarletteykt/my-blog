package auth

import (
	"github.com/scraletteykt/my-blog/internal/user"
)

type Options struct {
	Secret string
	Users  *user.Users
}
