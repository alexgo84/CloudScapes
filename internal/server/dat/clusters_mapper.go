package dat

import (
	"CloudScapes/pkg/wire"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type ClustersMapper struct {
	txn *sqlx.Tx
	ctx context.Context
}

type Cluster struct {
	ID      int64     `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created_at"`
	wire.NewCluster
}

func NewClustersMapper(ctx context.Context, txn *sqlx.Tx) ClustersMapper {
	return ClustersMapper{
		txn: txn,
		ctx: ctx,
	}
}

func (am *PlansMapper) CreateCluster(newCluster wire.NewCluster) (*Cluster, error) {
	c := Cluster{NewCluster: newCluster}
	err := namedGet(am.txn, "INSERT INTO clusters (name, accountid) VALUES (:name, :accountid) RETURNING id, created_at", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (am *PlansMapper) GetClusters(accountID int64) ([]Cluster, error) {
	var clusters []Cluster
	err := am.txn.SelectContext(am.ctx, &clusters, "select * from clusters WHERE accountid = $1 ORDER BY id desc", accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Cluster{}, nil
	}
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
