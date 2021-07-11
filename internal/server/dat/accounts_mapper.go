package dat

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type AccountsMapper struct {
	txn pgx.Tx
}

type Account struct {
	Id      int64     `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created_at"`
}

func NewAccountsMapper(txn pgx.Tx) AccountsMapper {
	return AccountsMapper{
		txn: txn,
	}
}

func (am *AccountsMapper) CreateAccount() (*Account, error) {
	return nil, nil
}

func (am *AccountsMapper) GetAccounts() ([]Account, error) {
	fmt.Printf("am.txn: %v\n", am.txn)

	// rows, err := am.txn.Query(context.Background(), "select * from accounts")
	// if err != nil {
	// 	return nil, err
	// }

	// var acc []Account
	// rows.Next()
	// if err := rows.Scan(&acc); err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
