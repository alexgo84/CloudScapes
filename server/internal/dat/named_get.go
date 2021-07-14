package dat

import "github.com/jmoiron/sqlx"

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
	return stmt.Get(arg, arg)
}
