package shared

import (
	"time"
)

type Deployment struct {
	ID         int64      `json:"id" db:"id"`
	Created    time.Time  `json:"created" db:"created_at"`
	CreatedBy  int64      `json:"createdBy" db:"created_by"`
	Modified   *time.Time `json:"modified" db:"modified_at"`
	ModifiedBy *int64     `json:"modifiedBy" db:"modified_by"`
	Deleted    *time.Time `json:"deleted" db:"deleted_at"`
	DeletedBy  *int64     `json:"deletedBy" db:"deleted_by"`

	Name               string  `json:"name" db:"name"`
	ImagePath          string  `json:"imagePath" db:"image_path"`
	ExcludeFromUpdates string  `json:"excludeFromUpdates" db:"exclude_from_updates"`
	PlanID             int64   `json:"planID" db:"planid"`
	SalesforceState    *string `json:"salesforce_state" db:"sf_state"`
	Replicas           *int64  `json:"replicas" db:"replicas"`

	CPULimit *int64 `json:"CPULimit" db:"cpu_limit"`
	MemLimit *int64 `json:"memLimit" db:"mem_limit"`
	CPUReq   *int64 `json:"CPUReq" db:"cpu_req"`
	MemReq   *int64 `json:"memReq" db:"mem_req"`

	DatabaseServiceName  int64 `json:"databaseServiceName" db:"database_service_name"`
	DatabaseServiceCloud int64 `json:"databaseServiceCloud" db:"database_service_cloud"`
	DatabaseServicePlan  int64 `json:"databaseServicePlan" db:"database_service_plan"`

	ClusterName string                 `json:"clusterName" db:"cluster_name"`
	EnvVars     map[string]interface{} `json:"envVars" db:"env_vars"`
	ConfigMaps  interface{}            `json:"configMaps" db:"config_maps"`
	CronJobs    interface{}            `json:"cronJobs" db:"cron_jobs"`
}
