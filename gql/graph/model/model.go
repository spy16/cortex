package model

import "github.com/chunked-app/cortex/block"

func ChunkModelFrom(c block.Chunk) *Chunk {
	res := &Chunk{
		ID:          c.ID,
		AuthorID:    c.Author,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Content:     c.Content,
		ContentType: ContentType(c.Type),
	}

	if c.Parent != "" {
		res.Parent = &c.Parent
	}

	if c.NextSibling != "" {
		res.NextSibling = &c.NextSibling
	}

	if c.PrevSibling != "" {
		res.PrevSibling = &c.PrevSibling
	}

	return res
}
