package database

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsUniqueConstraint(err error) bool {
	var pgerr *pgconn.PgError
	return errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation
}

func AsUniqueConstraint(err error) (*pgconn.PgError, bool) {
	var pgerr *pgconn.PgError
	return pgerr, errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation
}
