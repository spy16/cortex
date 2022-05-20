package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/chunked-app/cortex/gql/graph/model"
)

func (r *chunkMutationResolver) CreateChunk(ctx context.Context, request model.CreateRequest) (*model.Chunk, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chunkMutationResolver) UpdateChunk(ctx context.Context, id string, request model.UpdateRequest) (*model.Chunk, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chunkMutationResolver) DeleteChunk(ctx context.Context, id string) (*model.Chunk, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chunkQueryResolver) Chunk(ctx context.Context, id string) (*model.Chunk, error) {
	panic(fmt.Errorf("not implemented"))
}

// ChunkMutation returns ChunkMutationResolver implementation.
func (r *Resolver) ChunkMutation() ChunkMutationResolver { return &chunkMutationResolver{r} }

// ChunkQuery returns ChunkQueryResolver implementation.
func (r *Resolver) ChunkQuery() ChunkQueryResolver { return &chunkQueryResolver{r} }

type chunkMutationResolver struct{ *Resolver }
type chunkQueryResolver struct{ *Resolver }
