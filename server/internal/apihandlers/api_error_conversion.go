package apihandlers

import (
	"CloudScapes/pkg/wire"
	"database/sql"
	"fmt"
)

func convertToAPIIfNeeded(entity string, id int64, err error) error {
	if err == sql.ErrNoRows {
		return wire.NewNotFoundError(fmt.Sprintf("%s with id %d was not found", entity, id), err)
	}
	return err
}
