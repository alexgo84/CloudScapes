package wire

type NewPlan struct {
	AccountID int64 `json:"accountId" db:"accountid"`

	Name      string `json:"name" db:"name"`
	Replicas  int64  `json:"replicas" db:"replicas"`
	ClusterID int64  `json:"clusterID" db:"clusterid"`

	CPULimit string `json:"CPULimit" db:"cpu_limit"`
	MemLimit string `json:"memLimit" db:"mem_limit"`
	CPUReq   string `json:"CPUReq" db:"cpu_req"`
	MemReq   string `json:"memReq" db:"mem_req"`

	DatabaseServiceName  string `json:"databaseServiceName" db:"database_service_name"`
	DatabaseServiceCloud string `json:"databaseServiceCloud" db:"database_service_cloud"`
	DatabaseServicePlan  string `json:"databaseServicePlan" db:"database_service_plan"`

	EnvVars    map[string]interface{} `json:"envVars" db:"env_vars"`
	CronJobs   []CronJob              `json:"cronJobs" db:"cron_jobs"`
	ConfigMaps []ConfigMap            `json:"configMaps" db:"config_maps"`
}
