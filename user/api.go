package user

import (
	"context"
	"strings"
)

type Store interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, u User) error
}

type API struct{ Store Store }

func (api *API) User(ctx context.Context, id string) (*User, error) {
	id = strings.TrimSpace(id)

	return api.Store.GetUser(ctx, id)
}

func (api *API) Register(ctx context.Context, u User) (*User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := api.Store.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &u, nil
}
