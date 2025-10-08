package app

import "context"

type contextKeyAuth struct{}

func GetContextKeyAuth(ctx context.Context) (bool, bool) {
	auth, ok := ctx.Value(contextKeyAuth{}).(bool)
	return auth, ok
}

func SetContextKeyAuth(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, contextKeyAuth{}, v)
}
