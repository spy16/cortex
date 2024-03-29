package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/core/user"
	"github.com/chunked-app/cortex/server/gql/model"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, req *model.RegisterUserRequest) (*model.User, error) {
	u := user.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}

	regUser, err := r.UsersAPI.RegisterUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        regUser.ID,
		Name:      regUser.Name,
		Email:     regUser.Email,
		CreatedAt: regUser.CreatedAt,
		UpdatedAt: regUser.UpdatedAt,
	}, nil
}

func (r *mutationResolver) CreateChunk(ctx context.Context, req model.CreateRequest) (*model.Chunk, error) {
	curUser := user.From(ctx)

	ch := chunk.Chunk{
		Kind:   req.Kind,
		Data:   req.Data,
		Tags:   req.Tags,
		Author: curUser.ID,
	}

	createdCh, err := r.ChunksAPI.Create(ctx, ch)
	if err != nil {
		return nil, err
	}

	return model.ChunkFrom(*createdCh)
}

func (r *mutationResolver) UpdateChunk(ctx context.Context, id string, req model.UpdateRequest) (*model.Chunk, error) {
	upd, err := model.UpdatesFrom(req)
	if err != nil {
		return nil, err
	}

	updatedCh, err := r.ChunksAPI.Update(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	return model.ChunkFrom(*updatedCh)
}

func (r *mutationResolver) DeleteChunk(ctx context.Context, id string) (*model.Chunk, error) {
	ch, err := r.ChunksAPI.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkFrom(*ch)
}

func (r *queryResolver) Chunk(ctx context.Context, id string) (*model.Chunk, error) {
	ch, err := r.ChunksAPI.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkFrom(*ch)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	u, err := r.UsersAPI.User(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *userResolver) Chunks(ctx context.Context, obj *model.User, filter *model.ListFilter) ([]*model.Chunk, error) {
	opts := chunk.ListOptions{
		Author: obj.ID,
	}
	if filter != nil {
		if filter.Kind != nil {
			opts.Kind = *filter.Kind
		}
	}

	chunks, err := r.ChunksAPI.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var res []*model.Chunk
	for _, ch := range chunks {
		m, err := model.ChunkFrom(ch)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
