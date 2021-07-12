package dat

import (
	"CloudScapes/pkg/wire"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlansMapper struct {
	txn *sqlx.Tx
	ctx context.Context
}

type Plan struct {
	ID      int64     `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created_at"`
	wire.NewPlan
}

func NewPlansMapper(ctx context.Context, txn *sqlx.Tx) PlansMapper {
	return PlansMapper{
		txn: txn,
		ctx: ctx,
	}
}

func (am *PlansMapper) CreatePlan(newPlan wire.NewPlan) (*Plan, error) {
	p := Plan{NewPlan: newPlan}
	err := namedGet(am.txn, "INSERT INTO plans (accountid, name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_plan, env_vars, cron_jobs, config_maps) VALUES (:accountid, :name, :replicas, :clusterid, :cpu_limit, :mem_limit, :cpu_req, :mem_req, :database_service_name, :database_service_cloud, :database_service_plan, :env_vars, :cron_jobs, :config_maps) RETURNING id, created_at", &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (am *PlansMapper) UpdatePlan(planID int64, newPlan wire.NewPlan) (*Plan, error) {
	p := Plan{NewPlan: newPlan}
	err := namedGet(am.txn, "UPDATE plans SET (name, replicas, clusterid, cpu_limit, mem_limit, cpu_req, mem_req, database_service_name, database_service_cloud, database_service_plan, env_vars, cron_jobs, config_maps) VALUES (:name, :replicas, :clusterid, :cpu_limit, :mem_limit, :cpu_req, :mem_req, :database_service_name, :database_service_cloud, :database_service_plan, :env_vars, :cron_jobs, :config_maps) WHERE id = "+fmt.Sprintf("%d", planID)+" RETURNING id, created_at", &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (am *PlansMapper) GetPlans(accountID int64) ([]Plan, error) {
	var plans []Plan
	err := am.txn.SelectContext(am.ctx, &plans, "select * from plans WHERE accountid = $1 ORDER BY id desc", accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []Plan{}, nil
	}
	if err != nil {
		return nil, err
	}
	return plans, nil
}
