package listener

import (
	"CloudScapes/agent/internal/deployer"
	"CloudScapes/pkg/logger"
	"CloudScapes/pkg/shared"
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
)

type Listener struct {
	ch       <-chan *redis.Message
	chanName string
}

func NewListener(ch <-chan *redis.Message, chanName string) (*Listener, error) {
	if ch == nil {
		return nil, errors.New("cannot create a listener because the redis channel is nil")
	}
	return &Listener{
		ch:       ch,
		chanName: chanName,
	}, nil
}

func (l *Listener) ListenAndServe() {

	for msg := range l.ch {
		payload := []byte(msg.Payload)
		if err := handleAgentRequest(payload); err != nil {
			logger.Log(logger.ERROR, "failed to handle request", logger.Str("channel", msg.Channel), logger.Str("message", msg.String()))
		}
	}
}

func handleAgentRequest(req []byte) error {
	var agentReq shared.AgentRequest
	if err := json.Unmarshal(req, &agentReq); err != nil {
		return err
	}

	if agentReq.Deploy != nil {
		logger.Log(logger.INFO, "serving request", logger.Str("type", string(agentReq.Type)), logger.Str("name", agentReq.Deploy.Name), logger.Bool("async", agentReq.AsyncMode))
	} else {
		logger.Log(logger.INFO, "serving request", logger.Str("type", string(agentReq.Type)), logger.Bool("async", agentReq.AsyncMode))
	}

	switch agentReq.Type {
	case shared.AgentRequestTypeDeploy:

		if err := deployCustomer(agentReq.Deploy, agentReq.AsyncMode); err != nil {
			return err
		}
	}
	return nil
}

func deployCustomer(deployment *shared.K8sDeployment, asyncMode bool) error {
	ctx := context.Background()
	d, err := deployer.NewDeployer(asyncMode)
	if err != nil {
		return err
	}

	return d.ApplySpec(ctx, deployment)
}
