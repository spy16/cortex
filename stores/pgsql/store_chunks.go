package pgsql

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/pkg/errors"
)

func (st *Store) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	var rec chunkRecord

	found, err := st.db.From("chunks").Where(goqu.Ex{"id": id}).ScanStruct(&rec)
	if !found {
		return nil, errors.ErrNotFound
	} else if err != nil {
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}

	// TODO: get tags.

	return rec.toChunk(), nil
}

func (st *Store) List(ctx context.Context, opts chunk.ListOptions) ([]chunk.Chunk, error) {
	filter := goqu.Ex{}
	if opts.Author != "" {
		filter["author"] = opts.Author
	}

	if opts.Kind != "" {
		filter["kind"] = opts.Kind
	}

	var records []chunkRecord
	if err := st.db.From("chunks").Where(filter).ScanStructs(&records); err != nil {
		return nil, err
	}

	var res []chunk.Chunk
	for _, rec := range records {
		c := rec.toChunk()
		res = append(res, *c)
	}
	return res, nil
}

func (st *Store) Create(ctx context.Context, c chunk.Chunk) error {
	// TODO implement me
	panic("implement me")
}

func (st *Store) Update(ctx context.Context, id string, upd chunk.Updates) (*chunk.Chunk, error) {
	// TODO implement me
	panic("implement me")
}

func (st *Store) Delete(ctx context.Context, id string) (*chunk.Chunk, error) {
	// TODO implement me
	panic("implement me")
}

type chunkRecord struct {
	SeqID     int64     `db:"seq_id"`
	ID        string    `db:"id"`
	Kind      string    `db:"kind"`
	Data      string    `db:"data"`
	Author    string    `db:"author"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (rec chunkRecord) toChunk() *chunk.Chunk {
	return &chunk.Chunk{
		ID:        rec.ID,
		Kind:      rec.Kind,
		Data:      rec.Data,
		Tags:      nil,
		Author:    rec.Author,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

func chunkRecordFrom(ch chunk.Chunk) chunkRecord {
	return chunkRecord{
		ID:        ch.ID,
		Kind:      ch.Kind,
		Data:      ch.Data,
		Author:    ch.Author,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
	}
}
