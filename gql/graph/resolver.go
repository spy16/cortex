package graph

import (
	"context"

	"github.com/chunked-app/cortex/block"
	"github.com/chunked-app/cortex/user"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ChunksAPI ChunksAPI
	UsersAPI  UsersAPI
}

type UsersAPI interface {
	User(ctx context.Context, id string) (*user.User, error)
	Register(ctx context.Context, u user.User) (*user.User, error)
}

type ChunksAPI interface {
	Get(ctx context.Context, id string) (*block.Chunk, error)
	List(ctx context.Context, filter block.ListOptions) ([]block.Chunk, error)
	Create(ctx context.Context, c block.Chunk) (*block.Chunk, error)
	Update(ctx context.Context, id string, updates block.Updates) (*block.Chunk, error)
	Delete(ctx context.Context, id string) (*block.Chunk, error)
}
