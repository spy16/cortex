package user

import (
	"context"
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

type Store interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, u User) error
}

type API struct{ Store Store }

func (api *API) User(ctx context.Context, id string) (*User, error) {
	id = strings.TrimSpace(id)

	u, err := api.Store.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.
				WithMsgf("user with id '%s' not found", id).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	return u, nil
}

func (api *API) Register(ctx context.Context, u User) (*User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := api.Store.CreateUser(ctx, u); err != nil {
		if errors.Is(err, errors.ErrConflict) {
			return nil, errors.ErrConflict.
				WithMsgf("user with id '%s' already exists", u.ID).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}
	return &u, nil
}
