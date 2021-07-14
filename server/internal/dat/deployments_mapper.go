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
	txn       *sqlx.Tx
	ctx       context.Context
	accountID int64
}

type Deployment struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"accountId" db:"accountid"`
	Created   time.Time `json:"created" db:"created_at"`
	CreatedBy int64     `json:"createdBy" db:"created_by"`

	Modified   *time.Time `json:"modified" db:"modfied_at"`
	ModifiedBy *int64     `json:"modifiedBy" db:"modified_by"`
	Deleted    *time.Time `json:"deleted" db:"deleted_at"`
	DeletedBy  *int64     `json:"deletedBy" db:"deleted_by"`

	SalesforceState *string `json:"salesforceState" db:"salesforce_state"`

	wire.NewDeployment
}

func NewDeploymentsMapper(ctx context.Context, txn *sqlx.Tx, accountID int64) DeploymentsMapper {
	return DeploymentsMapper{
		txn:       txn,
		ctx:       ctx,
		accountID: accountID,
	}
}
func (am *DeploymentsMapper) CreateDeployment(newDeployment wire.NewDeployment, userID int64) (*Deployment, error) {
	d := Deployment{
		NewDeployment: newDeployment,
		AccountID:     am.accountID,
		CreatedBy:     userID,
	}
	err := namedGet(am.txn, `INSERT INTO deployments 
	(accountid, name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_deployment, env_vars, cron_jobs, config_maps, image_path, image_sha, exclude_from_updates, planid, salesforce_state) 
	VALUES 
	(:accountid, :name, :replicas, :clusterid, :cpu_limit, :mem_limit, :cpu_req, :mem_req, :database_service_name, :database_service_cloud, :database_service_deployment, :env_vars, :cron_jobs, :config_maps, :image_path, :image_sha, :exclude_from_updates, :planid, :salesforce_state) 
	RETURNING id, created_at, accountid`, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (am *DeploymentsMapper) UpdateDeployment(deploymentID, userID int64, newDeployment wire.NewDeployment) (*Deployment, error) {
	d := Deployment{
		NewDeployment: newDeployment,
		AccountID:     am.accountID,
		ModifiedBy:    &userID,
		ID:            deploymentID,
	}
	err := namedGet(am.txn, `UPDATE deployments SET  -- TODO: make it update for real
	name=:name, replicas=:replicas -- accountid, name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_deployment, env_vars, cron_jobs, config_maps, image_path, image_sha, exclude_from_updates, planid, salesforce_state
	WHERE id = :id
	RETURNING id, created_at, modified_at, accountid`, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (am *DeploymentsMapper) GetDeployments() ([]Deployment, error) {
	deployments := []Deployment{} // assign to empty array so that no result case does not return null
	err := am.txn.SelectContext(am.ctx, &deployments, "SELECT * FROM deployments WHERE accountid = $1 ORDER BY id desc", am.accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Deployment{}, nil
	}
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (am *DeploymentsMapper) GetDeployment(deploymentID int64) (*Deployment, error) {
	var deployment Deployment
	err := am.txn.GetContext(am.ctx, &deployment, "SELECT * FROM deployments WHERE accountid = $1 AND id = $2", am.accountID, deploymentID)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (am *DeploymentsMapper) DeleteDeployment(deploymentID int64) error {
	if _, err := am.GetDeployment(deploymentID); err != nil {
		return err
	}
	_, err := am.txn.QueryContext(am.ctx, "DELETE FROM deployments WHERE accountid = $1 AND id = $2", am.accountID, deploymentID)
	if err != nil {
		return err
	}
	return nil
}
