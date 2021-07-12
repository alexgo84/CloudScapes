package dat

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DataContext struct {
	Accounts AccountsMapper
}

func NewDataContext(ctx context.Context, txn *sqlx.Tx) DataContext {
	return DataContext{
		Accounts: NewAccountsMapper(ctx, txn),
	}
}
