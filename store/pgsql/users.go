package pgsql

import (
	"context"
	"database/sql"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/chunked-app/cortex/pkg/errors"
	"github.com/chunked-app/cortex/user"
)

func (st *PostgresQL) GetUser(ctx context.Context, id string) (*user.User, error) {
	return getUser(ctx, st.db, id)
}

func (st *PostgresQL) CreateUser(ctx context.Context, u user.User) error {
	q := sq.Insert(tableUsers).
		Columns("id", "name", "created_at", "updated_at").
		Values(u.ID, u.Name, u.CreatedAt, u.UpdatedAt).
		PlaceholderFormat(sq.Dollar)

	_, err := q.RunWith(st.db).ExecContext(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return errors.ErrConflict
		}
		return err
	}
	return nil
}

func getUser(ctx context.Context, r sq.BaseRunner, id string) (*user.User, error) {
	q := sq.Select("*").
		From(tableUsers).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	var rec user.User
	if err := q.RunWith(r).QueryRowContext(ctx).Scan(&rec.ID, &rec.Name, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound.WithMsgf("user with id '%s' not found", id)
		}
		return nil, err
	}

	return &rec, nil
}
