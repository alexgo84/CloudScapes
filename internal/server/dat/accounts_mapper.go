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
	Id          int64     `json:"id" db:"id"`
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
	var accounts []Account
	err := am.txn.SelectContext(am.ctx, &accounts, "select * from accounts")
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
