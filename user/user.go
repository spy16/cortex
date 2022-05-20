package user

import (
	"context"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	u.ID = strings.TrimSpace(u.ID)
	u.Name = strings.TrimSpace(u.Name)

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
		u.UpdatedAt = u.CreatedAt
	}

	if u.ID == "" {
		return errors.ErrInvalid.WithMsgf("user id must not be empty")
	}

	if u.Name == "" {
		return errors.ErrInvalid.WithMsgf("user name must not be empty")
	}

	return nil
}

func (u *User) IsEmpty() bool { return u.ID == "" }

func With(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func From(ctx context.Context) User {
	u, _ := ctx.Value(userKey).(User)
	return u
}

type key string

var userKey = key("user")
