package listener

import (
	"CloudScapes/pkg/logger"
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

func (l *Listener) ListenAndServe() error {

	for msg := range l.ch {
		payload := []byte(msg.Payload)
		if err := handleAgentRequest(payload); err != nil {
			logger.Log(logger.ERROR, "failed to handle request", logger.Str("channel", msg.Channel), logger.Str("message", msg.String()))
		}
	}

	return errors.New("should never go out of the loop")
}

func handleAgentRequest(req []byte) error {
	return nil
}
