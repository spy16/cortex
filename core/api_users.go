package core

import (
	"context"
	"fmt"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/core/user"
	"github.com/chunked-app/cortex/pkg/errors"
)

func (api *API) User(ctx context.Context, id string) (*user.User, error) {
	u, err := api.Users.FetchUser(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.
				WithMsgf("user with id '%s' not found", id).
				WithCausef(err.Error())
		}
		return nil, err
	}
	return u, nil
}

func (api *API) RegisterUser(ctx context.Context, u user.User) (*user.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := api.Users.CreateUser(ctx, u); err != nil {
		if errors.Is(err, errors.ErrConflict) {
			return nil, errors.ErrConflict.
				WithMsgf("user with id '%s' already exists", u.ID).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	ch := chunk.Chunk{
		ID:     fmt.Sprintf("u-%s", u.ID),
		Data:   &chunk.UserData{},
		Author: u.ID,
	}

	if _, err := api.createAny(ctx, ch); err != nil {
		return nil, err
	}

	return &u, nil
}
