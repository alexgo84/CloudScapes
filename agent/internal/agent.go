package agent

import (
	"CloudScapes/agent/internal/listener"
	"CloudScapes/pkg/pubsub"
	"os"
)

func Run() error {

	redisClient, err := pubsub.NewPubSubClient(nil)
	if err != nil {
		return err
	}

	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		clusterName = "test_cluster"
	}

	l, err := listener.NewListener(redisClient, clusterName)
	if err != nil {
		return err
	}

	l.ListenAndServe()
	return nil // should never reach here
}
