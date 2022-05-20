package chunk

import (
	"context"
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

// API provides functions to manage chunks.
type API struct{ Store Store }

func (api *API) Get(ctx context.Context, id string) (*Chunk, error) {
	id = strings.TrimSpace(id)
	if !idPattern.MatchString(id) {
		return nil, errors.ErrNotFound.
			WithMsgf("no chunk with id '%s'", id).
			WithCausef("id '%s' does not match pattern '%s'", id, idPattern)
	}

	ch, err := api.Store.Get(ctx, id)
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

func (api *API) List(ctx context.Context, opts ListOptions) ([]Chunk, error) {
	return api.Store.List(ctx, opts)
}

func (api *API) Create(ctx context.Context, ch Chunk) (*Chunk, error) {
	ch.genID()
	if err := ch.Validate(); err != nil {
		return nil, err
	}

	if err := api.Store.Create(ctx, ch); err != nil {
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}
	return &ch, nil
}

func (api *API) Update(ctx context.Context, id string, upd Updates) (*Chunk, error) {
	ch, err := api.Store.Update(ctx, id, upd)
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

func (api *API) Delete(ctx context.Context, id string) (*Chunk, error) {
	id = strings.TrimSpace(id)
	if !idPattern.MatchString(id) {
		return nil, errors.ErrNotFound.
			WithMsgf("no chunk with id '%s'", id).
			WithCausef("id '%s' does not match pattern '%s'", id, idPattern)
	}

	ch, err := api.Store.Delete(ctx, id)
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

type ListOptions struct {
	Parent string `json:"parent"`
}

func (opts ListOptions) IsMatch(ch Chunk) bool {
	parentMatch := opts.Parent == "" || ch.Parent == opts.Parent
	return parentMatch
}
