package dat

import (
	"CloudScapes/pkg/wire"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type DeploymentsMapper struct {
	txn *sqlx.Tx
	ctx context.Context
}

type Deployment struct {
	ID        int64     `json:"id" db:"id"`
	Created   time.Time `json:"created" db:"created_at"`
	CreatedBy int64     `json:"createdBy" db:"created_by"`

	Modified   *time.Time `json:"modified" db:"modfied_at"`
	ModifiedBy *int64     `json:"modifiedBy" db:"modified_by"`
	Deleted    *time.Time `json:"deleted" db:"deleted_at"`
	DeletedBy  *int64     `json:"deletedBy" db:"deleted_by"`

	SalesforceState *string `json:"salesforceState" db:"salesforce_state"`

	wire.NewDeployment
}

func NewDeploymentsMapper(ctx context.Context, txn *sqlx.Tx) DeploymentsMapper {
	return DeploymentsMapper{
		txn: txn,
		ctx: ctx,
	}
}

func (am *DeploymentsMapper) CreateDeployment(newDeployment wire.NewDeployment, userID int64, sfState *string) (*Deployment, error) {
	d := Deployment{NewDeployment: newDeployment, CreatedBy: userID, SalesforceState: sfState}
	err := namedGet(am.txn, `INSERT INTO deployments 
	(accountid, name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_plan, env_vars, cron_jobs, config_maps, salesforce_state) 
	VALUES 
	(:accountid, :name, :replicas, :clusterid, :cpu_limit, :mem_limit, :cpu_req, :mem_req, :database_service_name, :database_service_cloud, :database_service_plan, :env_vars, :cron_jobs, :config_maps, :salesforce_state) RETURNING id, created_at`, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (am *DeploymentsMapper) GetDeployment(accountID int64) ([]Deployment, error) {
	var deployments []Deployment
	err := am.txn.SelectContext(am.ctx, &deployments, "select * from plans WHERE accountid = $1 ORDER BY id desc", accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Deployment{}, nil
	}
	if err != nil {
		return nil, err
	}
	return deployments, nil
}
