package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/gql/graph/model"
)

func (r *chunkMutationResolver) CreateChunk(ctx context.Context, req model.CreateRequest) (*model.Chunk, error) {
	d, err := chunk.ParseData(req.Kind.String(), req.Data)
	if err != nil {
		return nil, err
	}

	ch := chunk.Chunk{
		Data:   d,
		Tags:   req.Tags,
		Author: req.AuthorID,
	}

	if req.Rank != nil {
		ch.Rank = *req.Rank
	}

	if req.ParentID != nil {
		ch.Parent = *req.ParentID
	}

	createdCh, err := r.ChunksAPI.Create(ctx, ch)
	if err != nil {
		return nil, err
	}

	return model.ChunkModelFrom(*createdCh)
}

func (r *chunkMutationResolver) UpdateChunk(ctx context.Context, id string, req model.UpdateRequest) (*model.Chunk, error) {
	upd, err := model.UpdatesFrom(req)
	if err != nil {
		return nil, err
	}

	updatedCh, err := r.ChunksAPI.Update(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	return model.ChunkModelFrom(*updatedCh)
}

func (r *chunkMutationResolver) DeleteChunk(ctx context.Context, id string) (*model.Chunk, error) {
	ch, err := r.ChunksAPI.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkModelFrom(*ch)
}

func (r *chunkQueryResolver) Chunk(ctx context.Context, id string) (*model.Chunk, error) {
	ch, err := r.ChunksAPI.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkModelFrom(*ch)
}

// ChunkMutation returns ChunkMutationResolver implementation.
func (r *Resolver) ChunkMutation() ChunkMutationResolver { return &chunkMutationResolver{r} }

// ChunkQuery returns ChunkQueryResolver implementation.
func (r *Resolver) ChunkQuery() ChunkQueryResolver { return &chunkQueryResolver{r} }

type chunkMutationResolver struct{ *Resolver }
type chunkQueryResolver struct{ *Resolver }
