package dat

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccountsMapper struct {
	txn *sqlx.Tx
	ctx context.Context
}

type Account struct {
	ID          int64     `json:"id" db:"id"`
	Created     time.Time `json:"created" db:"created_at"`
	CompanyName string    `json:"companyName" db:"company_name"`
}

func NewAccountsMapper(ctx context.Context, txn *sqlx.Tx) AccountsMapper {
	return AccountsMapper{
		txn: txn,
		ctx: ctx,
	}
}

func (am *AccountsMapper) CreateAccount(companyName string) (*Account, error) {
	acc := Account{CompanyName: companyName}
	if err := namedGet(am.txn, "INSERT INTO accounts (company_name) VALUES (:company_name) RETURNING id, created_at", &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (am *AccountsMapper) GetAccounts() ([]Account, error) {
	accounts := []Account{}
	err := am.txn.SelectContext(am.ctx, &accounts, "select * from accounts")
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

//GetLastAccount is a convenience method to provide account isolation before authentication
// is implemented by alwasy putting the latest account on context
func (am *AccountsMapper) GetLastAccount() (Account, error) {
	accounts := []Account{}
	err := am.txn.SelectContext(am.ctx, &accounts, "select * from accounts ORDER BY id desc LIMIT 1")
	if err != nil || len(accounts) == 0 {
		return Account{}, err
	}
	return accounts[0], nil
}
