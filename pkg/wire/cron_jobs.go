package wire

import (
	"encoding/json"
	"errors"
)

type CronJob struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	CPULimit        string            `json:"CPULimit"`
	MemLimit        string            `json:"memLimit"`
	CPUReq          string            `json:"CPUReq"`
	MemReq          string            `json:"memReq"`
	MinutesOverride *uint64           `json:"minutes"`
	HoursOverride   *uint64           `json:"hours"`
	Volumes         []Volume          `json:"volumes"`
	Args            []string          `json:"arguments"`
	EnvVars         map[string]string `json:"envVars"`
}

type Volume struct {
	Name      string       `json:"name"`
	MountPath string       `json:"mountPath"`
	Source    VolumeSource `json:"source"`
}

type VolumeSource struct {
	SecretName     *string `json:"secretName"`
	LocalConfigMap *string `json:"localConfigMap"`
}

type CronJobs []CronJob

func (cj *CronJobs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &cj)
}
