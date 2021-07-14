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

func (am *ClustersMapper) CreateCluster(newCluster wire.NewCluster) (*Cluster, error) {
	c := Cluster{NewCluster: newCluster}
	err := namedGet(am.txn, "INSERT INTO clusters (name, accountid) VALUES (:name, :accountid) RETURNING id, created_at", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (am *ClustersMapper) GetClusters(accountID int64) ([]Cluster, error) {
	clusters := []Cluster{} // assign to empty array so that no result case does not return null
	err := am.txn.SelectContext(am.ctx, &clusters, "select * from clusters WHERE accountid = $1 ORDER BY id desc", accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Cluster{}, nil
	}
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (am *ClustersMapper) GetCluster(accountID, clusterID int64) (*Cluster, error) {
	var cluster Cluster
	err := am.txn.GetContext(am.ctx, &cluster, "select * from clusters WHERE accountid = $1 AND id = $2", accountID, clusterID)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (am *ClustersMapper) DeleteCluster(accountID, clusterID int64) error {
	if _, err := am.GetCluster(accountID, clusterID); err != nil {
		return err
	}
	_, err := am.txn.QueryContext(am.ctx, "DELETE from clusters WHERE accountid = $1 AND id = $2", accountID, clusterID)
	if err != nil {
		return err
	}
	return nil
}
