package user

import "context"

var userKey = key("user")

type key string

func Into(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func From(ctx context.Context) User {
	u, _ := ctx.Value(userKey).(User)
	return u
}
