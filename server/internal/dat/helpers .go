package dat

import (
	"CloudScapes/pkg/wire"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

type namedPreparer interface {
	PrepareNamed(string) (*sqlx.NamedStmt, error)
}

// namedGet is a work around to execute a names statement that also allows writing the return
// values into the named entity. This is mostly used to insert / update objects while reading
// the created default values
func namedGet(db namedPreparer, query string, arg interface{}) error {
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	if err := stmt.Get(arg, arg); err != nil {
		if msg, ok := isConstraintViolation(err); ok {
			return wire.NewConflictError(msg)
		}
		return err
	}
	return nil
}

func isConstraintViolation(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return pgErr.Detail, true
		}
	}
	return "", false
}
