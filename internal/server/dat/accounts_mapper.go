package dat

import (
	"github.com/jackc/pgx/v4"
)

type AccountsMapper struct {
	txn *pgx.Tx
}

type Account struct {
	id int64,
	created time.Time,
	wire.Account,
}
func NewAccountsMapper(txn *pgx.Tx) AccountsMapper {
	return AccountsMapper{
		txn: txn,
	}
}

type (am *AccountsMapper) CreateAccount(newAccount wire.NewAccount) (dat.Account, error) {

}