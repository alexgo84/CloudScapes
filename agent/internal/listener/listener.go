package listener

import (
	"CloudScapes/agent/internal/deployer"
	"CloudScapes/pkg/logger"
	"CloudScapes/pkg/pubsub"
	"CloudScapes/pkg/shared"
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
)

type Listener struct {
	ch       <-chan *redis.Message
	client   *pubsub.PubSubClient
	chanName string
}

func NewListener(redisClient *pubsub.PubSubClient, chanName string) (*Listener, error) {
	ch := redisClient.Subscribe(context.Background(), chanName)

	if ch == nil {
		return nil, errors.New("cannot create a listener because the redis channel is nil")
	}
	return &Listener{
		ch:       ch,
		client:   redisClient,
		chanName: chanName,
	}, nil
}

func (l *Listener) ListenAndServe() {
	for msg := range l.ch {
		payload := []byte(msg.Payload)
		var agentReq shared.AgentRequest
		if err := json.Unmarshal(payload, &agentReq); err != nil {
			l.publishErrorReadingRequest(msg.Payload, err)
			continue
		}

		if err := handleAgentRequest(agentReq); err != nil {
			l.publishErrorHandlingRequest(&agentReq, err)
		}
	}
}

func handleAgentRequest(agentReq shared.AgentRequest) error {

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

func (l *Listener) publishErrorHandlingRequest(msg *shared.AgentRequest, err error) {
	res := shared.AgentResponse{}
	if err := l.client.PublishResponse(context.Background(), shared.ControlChannelName, res); err != nil {
		logger.Log(logger.ERROR, "failed to write agent response on handle failure", logger.Str("cluster", l.chanName), logger.Err(err))
	}
}

func (l *Listener) publishErrorReadingRequest(msg string, err error) {
	res := shared.AgentResponse{}
	if err := l.client.PublishResponse(context.Background(), shared.ControlChannelName, res); err != nil {
		logger.Log(logger.ERROR, "failed to write agent response on reading request failure", logger.Str("cluster", l.chanName), logger.Err(err))
	}
}
