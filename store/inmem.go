package store

import (
	"context"
	"sync"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
	"github.com/chunked-app/cortex/user"
)

type InMemory struct {
	userMu sync.RWMutex
	users  map[string]user.User

	chunkMu sync.RWMutex
	chunks  map[string]chunk.Chunk
}

func (st *InMemory) GetUser(ctx context.Context, id string) (*user.User, error) {
	st.userMu.RLock()
	defer st.userMu.RUnlock()

	u, found := st.users[id]
	if !found {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (st *InMemory) CreateUser(ctx context.Context, u user.User) error {
	st.userMu.Lock()
	defer st.userMu.Unlock()

	if st.users == nil {
		st.users = map[string]user.User{}
	}

	if _, found := st.users[u.ID]; found {
		return errors.ErrConflict
	}
	st.users[u.ID] = u
	return nil
}

func (st *InMemory) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	st.chunkMu.RLock()
	defer st.chunkMu.RUnlock()

	c, found := st.chunks[id]
	if !found {
		return nil, errors.ErrNotFound
	}
	return &c, nil
}

func (st *InMemory) List(ctx context.Context, opts chunk.ListOptions) ([]chunk.Chunk, error) {
	st.chunkMu.RLock()
	defer st.chunkMu.RUnlock()

	var res []chunk.Chunk
	for _, c := range st.chunks {
		if opts.Parent != "" || c.Parent == opts.Parent {
			res = append(res, c)
		}
	}
	return res, nil
}

func (st *InMemory) Create(ctx context.Context, c chunk.Chunk) error {
	st.chunkMu.Lock()
	defer st.chunkMu.Unlock()

	if st.chunks == nil {
		st.chunks = map[string]chunk.Chunk{}
	}

	if _, exists := st.chunks[c.ID]; exists {
		return errors.ErrConflict
	}

	st.chunks[c.ID] = c
	return nil
}

func (st *InMemory) Update(ctx context.Context, id string, updates chunk.Updates) (*chunk.Chunk, error) {
	st.chunkMu.Lock()
	defer st.chunkMu.Unlock()

	c, found := st.chunks[id]
	if !found {
		return nil, errors.ErrNotFound
	}

	if updates.ContentType != "" {
		c.Type = updates.ContentType
	}

	if updates.Content != "" {
		c.Content = updates.Content
	}

	if updates.ParentID != "" {
		if _, exists := st.chunks[updates.ParentID]; updates.ParentID != "" && !exists {
			return nil, errors.ErrNotFound.WithMsgf("parent with id '%s' not found", updates.ParentID)
		}
		c.Parent = updates.ParentID
	}

	if updates.PrevSibling != "" {
		if _, exists := st.chunks[updates.ParentID]; updates.ParentID != "" && !exists {
			return nil, errors.ErrNotFound.WithMsgf("parent with id '%s' not found", updates.ParentID)
		}
		c.Parent = updates.ParentID
	}

	st.chunks[c.ID] = c
	return &c, nil
}

func (st *InMemory) Delete(ctx context.Context, id string) (*chunk.Chunk, error) {
	st.chunkMu.Lock()
	defer st.chunkMu.Unlock()

	c, ok := st.chunks[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	// chunk with children cannot be deleted.
	for _, child := range st.chunks {
		if child.Parent == id {
			return nil, errors.ErrInvalid
		}
	}

	delete(st.chunks, id)
	return &c, nil
}
