package dat

import (
	"CloudScapes/pkg/wire"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlansMapper struct {
	txn       *sqlx.Tx
	ctx       context.Context
	accountID int64
}

type Plan struct {
	ID        int64     `json:"id" db:"id"`
	Created   time.Time `json:"created" db:"created_at"`
	AccountID int64     `json:"accountId" db:"accountid"`
	wire.NewPlan
}

func NewPlansMapper(ctx context.Context, txn *sqlx.Tx, accountID int64) PlansMapper {
	return PlansMapper{
		txn:       txn,
		ctx:       ctx,
		accountID: accountID,
	}
}

func (am *PlansMapper) CreatePlan(newPlan wire.NewPlan) (*Plan, error) {
	p := Plan{NewPlan: newPlan, AccountID: am.accountID}
	err := namedGet(am.txn, `INSERT INTO plans 
	(accountid, name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_plan, env_vars, cron_jobs, config_maps) 
	VALUES 
	(:accountid, :name, :replicas, :clusterid, :cpu_limit, :mem_limit, :cpu_req, :mem_req, :database_service_name, :database_service_cloud, :database_service_plan, :env_vars, :cron_jobs, :config_maps) RETURNING id, created_at`, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (am *PlansMapper) UpdatePlan(planID int64, newPlan wire.NewPlan) (*Plan, error) {
	p := Plan{NewPlan: newPlan}
	p.ID = planID
	err := namedGet(am.txn, `UPDATE plans SET  -- TODO: make it update for real
	name=:name, replicas=:replicas -- replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_plan, env_vars, cron_jobs, config_maps) 
	WHERE id = :id
	RETURNING id, created_at, accountid, modified_at`, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (am *PlansMapper) GetPlans() ([]Plan, error) {
	plans := []Plan{} // assign to empty array so that no result case does not return null
	err := am.txn.SelectContext(am.ctx, &plans, "SELECT * FROM plans WHERE accountid = $1 ORDER BY id desc", am.accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Plan{}, nil
	}
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (am *PlansMapper) GetPlan(planID int64) (*Plan, error) {
	var plan Plan
	err := am.txn.GetContext(am.ctx, &plan, "SELECT * FROM plans WHERE accountid = $1 AND id = $2", am.accountID, planID)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (am *PlansMapper) DeletePlan(planID int64) error {
	if _, err := am.GetPlan(planID); err != nil {
		return err
	}
	_, err := am.txn.QueryContext(am.ctx, "DELETE FROM plans WHERE accountid = $1 AND id = $2", am.accountID, planID)
	if err != nil {
		return err
	}
	return nil
}
