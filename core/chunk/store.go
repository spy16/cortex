package chunk

import "context"

// Store implementation is responsible for persisting chunks.
type Store interface {
	Get(ctx context.Context, id string) (*Chunk, error)
	List(ctx context.Context, opts ListOptions) ([]Chunk, error)
	Create(ctx context.Context, c Chunk) error
	Update(ctx context.Context, id string, upd Updates) (*Chunk, error)
	Delete(ctx context.Context, id string) (*Chunk, error)
}

type ListOptions struct {
	Parent string `json:"parent"`
	Author string `json:"author"`
}

func (opts ListOptions) IsMatch(ch Chunk) bool {
	parentMatch := opts.Parent == "" || ch.Parent == opts.Parent
	authorMatch := opts.Author == "" || ch.Author == opts.Author
	return parentMatch && authorMatch
}
