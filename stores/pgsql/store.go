package pgsql

import (
	"database/sql"
	_ "embed"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // goqu dialect
	_ "github.com/lib/pq"                               // pg driver

	"github.com/doug-martin/goqu/v9"
)

//go:embed schema.sql
var schema string

func Open(conStr string) (*Store, error) {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}
	st := Store{db: goqu.New("postgres", db)}

	if err := st.init(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &st, nil
}

type Store struct {
	db *goqu.Database
}

func (st *Store) init() error {
	_, err := st.db.Exec(schema)
	return err
}
