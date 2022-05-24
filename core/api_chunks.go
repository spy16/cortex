package core

import (
	"context"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

func (api *API) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	ch, err := api.Chunks.Get(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.
				WithMsgf("no chunk with id '%s'", id).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	return ch, nil
}

func (api *API) Exists(ctx context.Context, id string) (bool, error) {
	_, err := api.Get(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (api *API) List(ctx context.Context, opts chunk.ListOptions) ([]chunk.Chunk, error) {
	return api.Chunks.List(ctx, opts)
}

func (api *API) Create(ctx context.Context, ch chunk.Chunk) (*chunk.Chunk, error) {
	if err := ch.Validate(); err != nil {
		return nil, err
	}

	if _, err := api.User(ctx, ch.Author); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrInvalid.WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	if err := api.Chunks.Create(ctx, ch); err != nil {
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}
	return &ch, nil
}

func (api *API) Update(ctx context.Context, id string, upd chunk.Updates) (*chunk.Chunk, error) {
	ch, err := api.Chunks.Update(ctx, id, upd)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.
				WithMsgf("no chunk with id '%s'", id).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}
	return ch, nil
}

func (api *API) Delete(ctx context.Context, id string) (*chunk.Chunk, error) {
	ch, err := api.Chunks.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.
				WithMsgf("no chunk with id '%s'", id).
				WithCausef(err.Error())
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	return ch, nil
}
