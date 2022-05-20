package inmem

import (
	"context"
	"sync"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
	"github.com/chunked-app/cortex/user"
)

type Store struct {
	cMu    sync.RWMutex
	chunks map[string]chunk.Chunk

	uMu   sync.RWMutex
	users map[string]user.User
}

func (mem *Store) GetUser(ctx context.Context, id string) (*user.User, error) {
	mem.uMu.RLock()
	defer mem.uMu.RUnlock()

	u, found := mem.users[id]
	if !found {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (mem *Store) CreateUser(ctx context.Context, u user.User) error {
	mem.uMu.Lock()
	defer mem.uMu.Unlock()

	if mem.users == nil {
		mem.users = map[string]user.User{}
	} else if _, exists := mem.users[u.ID]; exists {
		return errors.ErrConflict
	}

	mem.users[u.ID] = u
	return nil
}

func (mem *Store) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	mem.cMu.RLock()
	defer mem.cMu.RUnlock()

	c, found := mem.chunks[id]
	if !found {
		return nil, errors.ErrNotFound
	}

	return &c, nil
}

func (mem *Store) List(ctx context.Context, opts chunk.ListOptions) ([]chunk.Chunk, error) {
	mem.cMu.RLock()
	defer mem.cMu.RUnlock()

	var res []chunk.Chunk
	for _, ch := range mem.chunks {
		if opts.IsMatch(ch) {
			res = append(res, ch)
		}
	}
	return res, nil
}

func (mem *Store) Create(ctx context.Context, c chunk.Chunk) error {
	mem.cMu.Lock()
	defer mem.cMu.Unlock()

	if mem.chunks == nil {
		mem.chunks = map[string]chunk.Chunk{}
	}

	if _, exists := mem.chunks[c.ID]; exists {
		return errors.ErrConflict
	}

	mem.chunks[c.ID] = c
	return nil
}

func (mem *Store) Update(ctx context.Context, id string, upd chunk.Updates) (*chunk.Chunk, error) {
	mem.cMu.Lock()
	defer mem.cMu.Unlock()

	if mem.chunks != nil {
		mem.chunks = map[string]chunk.Chunk{}
	}

	c, exists := mem.chunks[id]
	if !exists {
		return nil, errors.ErrNotFound
	}
	c.Apply(upd)

	mem.chunks[id] = c
	return &c, nil
}

func (mem *Store) Delete(ctx context.Context, id string) (*chunk.Chunk, error) {
	mem.cMu.Lock()
	defer mem.cMu.Unlock()

	ch, found := mem.chunks[id]
	if !found {
		return nil, errors.ErrNotFound
	}
	delete(mem.chunks, id)
	return &ch, nil
}
