package shared

type AgentRequestType string

const (
	AgentRequestTypeDeploy AgentRequestType = "deploy"
)

type AgentRequest struct {
	Type      AgentRequestType `json:"type"`
	Deploy    *K8sDeployment   `json:"deploy"`
	AsyncMode bool             `json:"asyncMode"`
}
