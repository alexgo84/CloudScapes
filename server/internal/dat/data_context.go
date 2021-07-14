package dat

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DataContext struct {
	Accounts    AccountsMapper
	Users       UsersMapper
	Clusters    ClustersMapper
	Plans       PlansMapper
	Deployments DeploymentsMapper
}

func NewDataContext(ctx context.Context, txn *sqlx.Tx, accountID int64) DataContext {
	return DataContext{
		Accounts:    NewAccountsMapper(ctx, txn),
		Users:       NewUsersMapper(ctx, txn, accountID),
		Clusters:    NewClustersMapper(ctx, txn, accountID),
		Plans:       NewPlansMapper(ctx, txn, accountID),
		Deployments: NewDeploymentsMapper(ctx, txn, accountID),
	}
}
