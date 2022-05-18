package pgsql

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/lib/pq" // postgresql driver
)

//go:embed schema.sql
var schema string

const (
	tableChunks = "chunks"
	tableUsers  = "users"
)

func Open(spec string) (*PostgresQL, error) {
	db, err := sql.Open("postgres", spec)
	if err != nil {
		return nil, err
	}
	st := &PostgresQL{db: db}
	return st, st.init()
}

type PostgresQL struct{ db *sql.DB }

func (st *PostgresQL) init() error {
	_, err := st.db.Exec(schema)
	return err
}

func (st *PostgresQL) Close() error { return st.db.Close() }

func (st *PostgresQL) withTx(ctx context.Context, fns ...txFn) error {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, fn := range fns {
		if err := fn(ctx, tx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

type txFn func(ctx context.Context, tx *sql.Tx) error
