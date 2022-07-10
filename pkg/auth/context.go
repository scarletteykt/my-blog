package auth

import (
	"context"
)

var (
	contextKeyUser = contextKey("ctx_user")

	failedContextUser = User{
		ID:       -1,
		Username: "",
	}
)

func FromContext(ctx context.Context) User {
	value := ctx.Value(contextKeyUser)
	if value == nil {
		return failedContextUser
	}

	user, ok := value.(User)
	if !ok {
		return failedContextUser
	}

	return user
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}
