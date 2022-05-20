package inmem

import (
	"context"
	"sync"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]chunk.Chunk
}

func (mem *Store) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	c, found := mem.data[id]
	if !found {
		return nil, errors.ErrNotFound
	}

	return &c, nil
}

func (mem *Store) List(ctx context.Context, opts chunk.ListOptions) ([]chunk.Chunk, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	var res []chunk.Chunk
	for _, ch := range mem.data {
		if opts.IsMatch(ch) {
			res = append(res, ch)
		}
	}
	return res, nil
}

func (mem *Store) Create(ctx context.Context, c chunk.Chunk) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	if mem.data == nil {
		mem.data = map[string]chunk.Chunk{}
	}

	if _, exists := mem.data[c.ID]; exists {
		return errors.ErrConflict
	}

	mem.data[c.ID] = c
	return nil
}

func (mem *Store) Update(ctx context.Context, id string, upd chunk.Updates) (*chunk.Chunk, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	if mem.data != nil {
		mem.data = map[string]chunk.Chunk{}
	}

	c, exists := mem.data[id]
	if !exists {
		return nil, errors.ErrNotFound
	}
	c.Apply(upd)

	mem.data[c.ID] = c
	return &c, nil
}

func (mem *Store) Delete(ctx context.Context, id string) (*chunk.Chunk, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	ch, found := mem.data[id]
	if !found {
		return nil, errors.ErrNotFound
	}
	delete(mem.data, id)
	return &ch, nil
}
