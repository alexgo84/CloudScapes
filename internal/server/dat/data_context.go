package dat

import "github.com/jackc/pgx/v4"

type DataContext struct {
	Accounts AccountsMapper
}

func NewDataContext(txn *pgx.Tx) DataContext {
	return DataContext{
		Accounts: NewAccountsMapper(txn),
	}
}
