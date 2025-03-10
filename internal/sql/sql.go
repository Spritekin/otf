/*
Package sql implements persistent storage using the postgres database.
*/
package sql

import (
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/leg100/otf/internal"
)

// String converts a go-string into a postgres non-null string
func String(s string) pgtype.Text {
	return pgtype.Text{String: s, Status: pgtype.Present}
}

// NullString returns a postgres null string
func NullString() pgtype.Text {
	return pgtype.Text{Status: pgtype.Null}
}

// UUID converts a google-go-uuid into a postgres non-null UUID
func UUID(s uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: s, Status: pgtype.Present}
}

// Timestamptz converts a go-time into a postgres non-null timestamptz
func Timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Status: pgtype.Present}
}

func Error(err error) error {
	var pgErr *pgconn.PgError
	switch {
	case noRowsInResultError(err):
		return internal.ErrResourceNotFound
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case "23503": // foreign key violation
			return &internal.ForeignKeyError{PgError: pgErr}
		case "23505": // unique violation
			return internal.ErrResourceAlreadyExists
		}
		fallthrough
	default:
		return err
	}
}

func noRowsInResultError(err error) bool {
	for {
		err = errors.Unwrap(err)
		if err == nil {
			return false
		} else if err.Error() == "no rows in result set" {
			return true
		}
	}
}
