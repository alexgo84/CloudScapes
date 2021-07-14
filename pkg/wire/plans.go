package wire

import (
	"encoding/json"
	"errors"
)

type NewPlan struct {
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

	EnvVars    StringInterfaceMap `json:"envVars" db:"env_vars"`
	CronJobs   CronJobs           `json:"cronJobs" db:"cron_jobs"`
	ConfigMaps ConfigMaps         `json:"configMaps" db:"config_maps"`
}

type StringInterfaceMap map[string]interface{}

func (m *StringInterfaceMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}
