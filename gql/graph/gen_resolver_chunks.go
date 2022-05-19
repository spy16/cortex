package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/chunked-app/cortex/block"
	"github.com/chunked-app/cortex/gql/graph/model"
	"github.com/chunked-app/cortex/user"
)

func (r *chunkResolver) Author(ctx context.Context, obj *model.Chunk) (*model.User, error) {
	return r.ChunkQuery().User(ctx, obj.AuthorID)
}

func (r *chunkResolver) Children(ctx context.Context, obj *model.Chunk) ([]*model.Chunk, error) {
	log.Debugf("looking for chunks with parent=%s", obj.ID)

	children, err := r.ChunksAPI.List(ctx, block.ListOptions{Parent: obj.ID})
	if err != nil {
		return nil, err
	}

	var res []*model.Chunk
	for _, c := range children {
		log.Debugf("added %s to parent %s", c.ID, obj.ID)
		res = append(res, model.ChunkModelFrom(c))
	}
	return res, nil
}

func (r *chunkMutationResolver) RegisterUser(ctx context.Context, req model.UserRegistrationRequest) (*model.User, error) {
	u := user.User{
		ID:   req.ID,
		Name: req.Name,
	}

	regUser, err := r.UsersAPI.Register(ctx, u)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        regUser.ID,
		Name:      regUser.Name,
		CreatedAt: regUser.CreatedAt,
		UpdatedAt: regUser.UpdatedAt,
	}, nil
}

func (r *chunkMutationResolver) CreateChunk(ctx context.Context, request model.CreateRequest) (*model.Chunk, error) {
	if request.ContentType == nil {
		t := model.ContentTypeText
		request.ContentType = &t
	}

	ch := block.Chunk{
		ID:      request.ID,
		Author:  request.AuthorID,
		Type:    request.ContentType.String(),
		Content: request.Content,
	}
	if request.ParentID != nil {
		ch.Parent = *request.ParentID
	}
	if request.PrevSibling != nil {
		ch.PrevSibling = *request.PrevSibling
	}

	createdChunk, err := r.ChunksAPI.Create(ctx, ch)
	if err != nil {
		return nil, err
	}

	return model.ChunkModelFrom(*createdChunk), nil
}

func (r *chunkMutationResolver) UpdateChunk(ctx context.Context, id string, request model.UpdateRequest) (*model.Chunk, error) {
	upd := block.Updates{}
	if request.Content != nil {
		upd.Content = *request.Content
	}

	if request.ParentID != nil {
		upd.ParentID = *request.ParentID
	}

	if request.PrevSibling != nil {
		upd.PrevSibling = *request.PrevSibling
	}

	if request.ContentType != nil {
		t := request.ContentType.String()
		upd.ContentType = t
	}

	result, err := r.ChunksAPI.Update(ctx, id, upd)
	if err != nil {
		return nil, err
	}
	return model.ChunkModelFrom(*result), nil
}

func (r *chunkMutationResolver) DeleteChunk(ctx context.Context, id string) (*model.Chunk, error) {
	c, err := r.ChunksAPI.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkModelFrom(*c), nil
}

func (r *chunkQueryResolver) User(ctx context.Context, id string) (*model.User, error) {
	u, err := r.UsersAPI.User(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *chunkQueryResolver) Chunk(ctx context.Context, id string) (*model.Chunk, error) {
	c, err := r.ChunksAPI.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ChunkModelFrom(*c), nil
}

// Chunk returns ChunkResolver implementation.
func (r *Resolver) Chunk() ChunkResolver { return &chunkResolver{r} }

// ChunkMutation returns ChunkMutationResolver implementation.
func (r *Resolver) ChunkMutation() ChunkMutationResolver { return &chunkMutationResolver{r} }

// ChunkQuery returns ChunkQueryResolver implementation.
func (r *Resolver) ChunkQuery() ChunkQueryResolver { return &chunkQueryResolver{r} }

type chunkResolver struct{ *Resolver }
type chunkMutationResolver struct{ *Resolver }
type chunkQueryResolver struct{ *Resolver }
