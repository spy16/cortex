package chunk

import "context"

type API struct {
	Store interface {
		Create(ctx context.Context, c Chunk) error
	}
}

func (api *API) Create(ctx context.Context, c Chunk) (*Chunk, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := api.Store.Create(ctx, c); err != nil {
		return nil, err
	}

	return &c, nil
}
