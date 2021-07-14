package wire

type NewDeployment struct {
	Name string `json:"name" db:"name"`

	Replicas  *int64 `json:"replicas" db:"replicas"`
	ClusterID *int64 `json:"clusterID" db:"clusterid"`

	ImagePath          string `json:"imagePath" db:"image_path"`
	ImageSHA           string `json:"imageSHA" db:"image_sha"`
	ExcludeFromUpdates bool   `json:"excludeFromUpdates" db:"exclude_from_updates"`
	PlanID             int64  `json:"planID" db:"planid"`

	CPULimit *string `json:"CPULimit" db:"cpu_limit"`
	MemLimit *string `json:"memLimit" db:"mem_limit"`
	CPUReq   *string `json:"CPUReq" db:"cpu_req"`
	MemReq   *string `json:"memReq" db:"mem_req"`

	DatabaseServiceName  *string `json:"databaseServiceName" db:"database_service_name"`
	DatabaseServiceCloud *string `json:"databaseServiceCloud" db:"database_service_cloud"`
	DatabaseServicePlan  *string `json:"databaseServicePlan" db:"database_service_plan"`

	EnvVars    StringInterfaceMap `json:"envVars" db:"env_vars"`
	CronJobs   CronJobs           `json:"cronJobs" db:"cron_jobs"`
	ConfigMaps []ConfigMap        `json:"configMaps" db:"config_maps"`
}
