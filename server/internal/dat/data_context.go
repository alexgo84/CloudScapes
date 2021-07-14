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

func NewDataContext(ctx context.Context, txn *sqlx.Tx) DataContext {
	return DataContext{
		Accounts:    NewAccountsMapper(ctx, txn),
		Users:       NewUsersMapper(ctx, txn),
		Clusters:    NewClustersMapper(ctx, txn),
		Plans:       NewPlansMapper(ctx, txn),
		Deployments: NewDeploymentsMapper(ctx, txn),
	}
}
