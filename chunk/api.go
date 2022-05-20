package chunk

import "context"

type Store interface {
	Create(ctx context.Context, c Chunk) error
}

type API struct {
	Store Store
}

type ListOptions struct {
	Parent string `json:"parent"`
}
