package auth

import (
	"context"
)

var (
	contextKeyUser = contextKey("ctx_user")
)

func FromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(contextKeyUser).(User)
	return user, ok
}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}
