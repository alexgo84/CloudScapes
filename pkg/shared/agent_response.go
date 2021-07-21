package shared

const ControlChannelName = "control"

type AgentResponseType string

const (
	AgentResponseTypeOK    AgentResponseType = "ok"
	AgentResponseTypeError AgentResponseType = "error"
)

type AgentResponse struct {
	Type    AgentResponseType `json:"type"`
	Deploy  *K8sDeployment    `json:"deploy"`
	Message string            `json:"message"`
}
