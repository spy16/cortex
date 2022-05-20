package model

import (
	"encoding/json"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
	"github.com/chunked-app/cortex/user"
)

func UserFrom(u user.User) (*User, error) {
	return &User{
		ID:        u.ID,
		Name:      u.Name,
		Chunks:    nil,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func ChunkFrom(c chunk.Chunk) (*Chunk, error) {
	b, err := json.Marshal(c.Data)
	if err != nil {
		return nil, errors.ErrInternal.
			WithMsgf("failed to marshal chunk data").
			WithCausef(err.Error())
	}

	res := &Chunk{
		ID:        c.ID,
		Rank:      c.Rank,
		Kind:      c.Kind,
		Data:      string(b),
		Tags:      c.Tags,
		AuthorID:  c.Author,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.Parent != "" {
		res.ParentID = &c.Parent
	}

	return res, nil
}

func UpdatesFrom(request UpdateRequest) (chunk.Updates, error) {
	var upd chunk.Updates

	if request.Data != nil {
		if request.Kind == nil {
			return upd, errors.ErrInvalid.WithMsgf("when updating 'data', 'kind' must be specified")
		}

		d, err := chunk.ParseData(*request.Kind, *request.Data)
		if err != nil {
			return upd, err
		}
		upd.Data = d
	}

	if request.Parent != nil {
		upd.Parent = *request.Parent
	}

	if request.Rank != nil {
		upd.Rank = *request.Rank
	}

	if len(request.Tags) > 0 {
		upd.Tags = request.Tags
	}

	return upd, nil
}
