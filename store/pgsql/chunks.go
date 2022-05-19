package pgsql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/chunked-app/cortex/pkg/errors"
)

func (st *PostgresQL) Get(ctx context.Context, id string) (*block.Chunk, error) {
	return getChunk(ctx, st.db, id)
}

func (st *PostgresQL) List(ctx context.Context, opts block.ListOptions) ([]block.Chunk, error) {
	q := sq.Select("*").
		From(tableChunks).
		PlaceholderFormat(sq.Dollar)

	if opts.Parent != "" {
		q = q.Where(sq.Eq{"parent_id": opts.Parent})
	} else {
		q = q.Where(sq.Eq{"parent_id": nil})
	}

	rows, err := q.RunWith(st.db).QueryContext(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var res []block.Chunk
	for rows.Next() {
		var rec chunkRecord
		if err := rows.Scan(&rec.ID, &rec.Author, &rec.CreatedAt, &rec.UpdatedAt, &rec.ParentID,
			&rec.PrevSiblingID, &rec.NextSiblingID, &rec.Content, &rec.ContentType); err != nil {
			return nil, err
		}
		ch := rec.toChunk()
		res = append(res, *ch)
	}

	if rows.Err() != nil {
		return res, err
	}
	return res, nil
}

func (st *PostgresQL) Create(ctx context.Context, c block.Chunk) error {
	rec := newChunkRecord(c)

	ensureUser := txFn(func(ctx context.Context, tx *sql.Tx) error {
		_, err := getUser(ctx, tx, c.Author)
		if err != nil {
			return err
		}
		return nil
	})

	createChunk := txFn(func(ctx context.Context, tx *sql.Tx) error {
		q := sq.Insert(tableChunks).
			Columns("id", "author", "created_at", "updated_at",
				"parent_id", "prev_sibling_id", "next_sibling_id",
				"content_type", "content").
			Values(rec.ID, rec.Author, rec.CreatedAt, rec.UpdatedAt,
				rec.ParentID, rec.PrevSiblingID, rec.NextSiblingID,
				rec.ContentType, rec.Content).
			PlaceholderFormat(sq.Dollar)

		_, err := q.RunWith(tx).ExecContext(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				return errors.ErrConflict
			}
			return err
		}
		return nil
	})

	return st.withTx(ctx, ensureUser, createChunk)
}

func (st *PostgresQL) Update(ctx context.Context, id string, updates block.Updates) (*block.Chunk, error) {
	q := sq.Update(tableChunks).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now())

	if updates.Content != "" {
		q = q.Set("content", updates.Content)
	}

	if updates.ContentType != "" {
		q = q.Set("content_type", updates.ContentType)
	}

	if updates.ParentID != "" {
		q = q.Set("parent_id", updates.ParentID)
	}

	if updates.PrevSibling != "" {
		q = q.Set("prev_sibling_id", updates.PrevSibling)
	}

	_, err := q.RunWith(st.db).ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return getChunk(ctx, st.db, id)
}

func (st *PostgresQL) Delete(ctx context.Context, id string) (*block.Chunk, error) {
	var c *block.Chunk

	checkCount := func(ctx context.Context, tx *sql.Tx) error {
		q := sq.Select("count(*)").
			From(tableChunks).
			Where(sq.Eq{"parent_id": id}).
			PlaceholderFormat(sq.Dollar)

		var count int
		if err := q.RunWith(tx).QueryRowContext(ctx).Scan(&count); err != nil {
			return err
		} else if count > 0 {
			return errors.ErrInvalid.WithMsgf("chunk with children cannot be deleted")
		}
		return nil
	}

	getFn := func(ctx context.Context, tx *sql.Tx) error {
		res, err := getChunk(ctx, tx, id)
		if err != nil {
			return err
		}
		c = res
		return nil
	}

	delFn := func(ctx context.Context, tx *sql.Tx) error {
		q := sq.Delete(tableChunks).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)
		_, err := q.RunWith(tx).ExecContext(ctx)
		return err
	}

	if txErr := st.withTx(ctx, checkCount, getFn, delFn); txErr != nil {
		return nil, txErr
	}
	return c, nil
}

func getChunk(ctx context.Context, r sq.BaseRunner, id string) (*block.Chunk, error) {
	var rec chunkRecord
	q := sq.Select("id", "author", "created_at", "updated_at", "parent_id",
		"prev_sibling_id", "next_sibling_id", "content", "content_type").
		From(tableChunks).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	row := q.RunWith(r).QueryRowContext(ctx)
	if err := row.Scan(&rec.ID, &rec.Author, &rec.CreatedAt, &rec.UpdatedAt, &rec.ParentID,
		&rec.PrevSiblingID, &rec.NextSiblingID, &rec.Content, &rec.ContentType); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return rec.toChunk(), nil
}

func newChunkRecord(c block.Chunk) (rec chunkRecord) {
	return chunkRecord{
		ID:            c.ID,
		Author:        c.Author,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		ParentID:      chunkIDValue(c.Parent),
		PrevSiblingID: chunkIDValue(c.PrevSibling),
		NextSiblingID: chunkIDValue(c.NextSibling),
		Content:       c.Content,
		ContentType:   c.Type,
	}
}

type chunkRecord struct {
	ID        string    `db:"id"`
	Author    string    `db:"author"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	ParentID      sql.NullString `db:"parent_id"`
	PrevSiblingID sql.NullString `db:"prev_sibling_id"`
	NextSiblingID sql.NullString `db:"next_sibling_id"`

	Content     string `db:"content"`
	ContentType string `db:"content_type"`
}

func (cr chunkRecord) toChunk() *block.Chunk {
	return &block.Chunk{
		ID:        cr.ID,
		Author:    cr.Author,
		CreatedAt: cr.CreatedAt,
		UpdatedAt: cr.UpdatedAt,

		Type:    cr.ContentType,
		Content: cr.Content,

		Parent:      cr.ParentID.String,
		NextSibling: cr.NextSiblingID.String,
		PrevSibling: cr.PrevSiblingID.String,
	}
}

func chunkIDValue(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: v,
		Valid:  true,
	}
}
