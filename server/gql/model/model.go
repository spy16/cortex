package model

import (
	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

func ChunkFrom(c chunk.Chunk) (*Chunk, error) {
	res := &Chunk{
		ID:        c.ID,
		Kind:      c.Kind,
		Data:      c.Data,
		Tags:      c.Tags,
		AuthorID:  c.Author,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	return res, nil
}

func UpdatesFrom(request UpdateRequest) (chunk.Updates, error) {
	var upd chunk.Updates

	if len(request.Tags) > 0 {
		upd.Tags = request.Tags
	}

	if (request.Kind != nil && request.Data == nil) || (request.Kind == nil && request.Data != nil) {
		return upd, errors.ErrInvalid.WithMsgf("only one of kind & data is set, both must be set")
	} else {
		upd.Kind = *request.Kind
		upd.Data = *request.Data
	}

	return upd, nil
}
