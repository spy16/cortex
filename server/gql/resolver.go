package gql

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/core/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UsersAPI  UsersAPI
	ChunksAPI ChunksAPI
}

type UsersAPI interface {
	User(ctx context.Context, id string) (*user.User, error)
	RegisterUser(ctx context.Context, u user.User) (*user.User, error)
}

type ChunksAPI interface {
	Get(ctx context.Context, id string) (*chunk.Chunk, error)
	List(ctx context.Context, filter chunk.ListOptions) ([]chunk.Chunk, error)
	Create(ctx context.Context, c chunk.Chunk) (*chunk.Chunk, error)
	Update(ctx context.Context, id string, updates chunk.Updates) (*chunk.Chunk, error)
	Delete(ctx context.Context, id string) (*chunk.Chunk, error)
}
