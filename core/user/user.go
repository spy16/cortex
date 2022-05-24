package user

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

const maxIDLen = 32

var userIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type Store interface {
	FetchUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, u User) error
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
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

	if u.Name == "" {
		return errors.ErrInvalid.WithMsgf("user name must not be empty")
	}

	return validateUserID(u.ID)
}

func (u *User) IsEmpty() bool { return u.ID == "" }

func validateUserID(id string) error {
	if len(id) == 0 || len(id) > maxIDLen {
		return errors.ErrInvalid.WithMsgf("login_id must have 1-%d characters", maxIDLen)
	} else if !userIDPattern.MatchString(id) {
		return errors.ErrInvalid.WithMsgf("login_id must match pattern '%s'", userIDPattern)
	}
	return nil
}
