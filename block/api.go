package block

import (
	"context"
	"sort"

	"github.com/chunked-app/cortex/pkg/errors"
)

type Store interface {
	Get(ctx context.Context, id string) (*Chunk, error)
	List(ctx context.Context, opts ListOptions) ([]Chunk, error)
	Create(ctx context.Context, c Chunk) error
	Delete(ctx context.Context, id string) (*Chunk, error)
	Update(ctx context.Context, id string, updates Updates) (*Chunk, error)
}

type API struct{ Store Store }

// Get returns the chunk with given identifier.
func (api API) Get(ctx context.Context, id string) (*Chunk, error) {
	c, err := api.Store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.WithMsgf("chunk with id '%s' not found", id)
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	// TODO: add authorisation check here.
	return c, nil
}

func (api API) List(ctx context.Context, opts ListOptions) ([]Chunk, error) {
	chunks, err := api.Store.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	if opts.SiblingOrder {
		sort.Slice(chunks, func(i, j int) bool {
			return chunks[i].NextSibling > chunks[j].NextSibling
		})
	}

	return chunks, nil
}

func (api API) Create(ctx context.Context, chunk Chunk) (*Chunk, error) {
	if err := chunk.Validate(); err != nil {
		return nil, err
	}

	err := api.Store.Create(ctx, chunk)
	if err != nil {
		return nil, err
	}

	return &chunk, nil
}

func (api API) Update(ctx context.Context, id string, updates Updates) (*Chunk, error) {
	if err := validateID(id, true); err != nil {
		return nil, err
	}

	updated, err := api.Store.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (api API) Delete(ctx context.Context, id string) (*Chunk, error) {
	if err := validateID(id, true); err != nil {
		return nil, err
	}

	c, err := api.Store.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

type ListOptions struct {
	Parent       string
	SiblingOrder bool
}
