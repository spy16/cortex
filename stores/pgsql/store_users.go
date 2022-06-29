package pgsql

import (
	"context"

	"github.com/chunked-app/cortex/core/user"
)

func (st *Store) FetchUser(ctx context.Context, id string) (*user.User, error) {
	// TODO implement me
	panic("implement me")
}

func (st *Store) CreateUser(ctx context.Context, u user.User) error {
	// TODO implement me
	panic("implement me")
}
