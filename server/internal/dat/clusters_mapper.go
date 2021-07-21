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
	txn       *sqlx.Tx
	ctx       context.Context
	accountID int64
}

type Cluster struct {
	ID        int64     `json:"id" db:"id"`
	Created   time.Time `json:"created" db:"created_at"`
	AccountID int64     `json:"accountId" db:"accountid"`
	wire.NewCluster
}

func NewClustersMapper(ctx context.Context, txn *sqlx.Tx, accountID int64) ClustersMapper {
	return ClustersMapper{
		txn:       txn,
		ctx:       ctx,
		accountID: accountID,
	}
}

func (am *ClustersMapper) CreateCluster(newCluster wire.NewCluster) (*Cluster, error) {
	c := Cluster{NewCluster: newCluster, AccountID: am.accountID}
	err := namedGet(am.txn, "INSERT INTO clusters (name, accountid) VALUES (:name, :accountid) RETURNING id, created_at", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (am *ClustersMapper) GetClusters() ([]Cluster, error) {
	clusters := []Cluster{} // assign to empty array so that no result case does not return null
	err := am.txn.SelectContext(am.ctx, &clusters, "select * from clusters WHERE accountid = $1 ORDER BY id desc", am.accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Cluster{}, nil
	}
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (am *ClustersMapper) GetCluster(clusterID int64) (*Cluster, error) {
	var cluster Cluster
	err := am.txn.GetContext(am.ctx, &cluster, "select * from clusters WHERE accountid = $1 AND id = $2", am.accountID, clusterID)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (am *ClustersMapper) DeleteCluster(clusterID int64) error {
	if _, err := am.GetCluster(clusterID); err != nil {
		return err
	}
	_, err := am.txn.QueryContext(am.ctx, "DELETE from clusters WHERE accountid = $1 AND id = $2", am.accountID, clusterID)
	if err != nil {
		if msg, ok := isConstraintViolation(err); ok {
			return wire.NewConflictError(msg, err)
		}
		return err
	}
	return nil
}
