package model

import (
	"encoding/json"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

func ChunkModelFrom(c chunk.Chunk) (*Chunk, error) {
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
