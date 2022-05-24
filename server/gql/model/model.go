package model

import (
	"encoding/json"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

func ChunkFrom(c chunk.Chunk) (*Chunk, error) {
	b, err := json.Marshal(c.Data)
	if err != nil {
		return nil, errors.ErrInternal.
			WithMsgf("failed to marshal chunk data").
			WithCausef(err.Error())
	}

	res := &Chunk{
		ID:        c.ID,
		Kind:      c.Data.Kind(),
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

	if len(request.Tags) > 0 {
		upd.Tags = request.Tags
	}

	return upd, nil
}
