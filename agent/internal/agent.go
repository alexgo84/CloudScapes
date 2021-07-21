package agent

import (
	"CloudScapes/agent/internal/listener"
	"CloudScapes/pkg/redis"
	"context"
	"os"
)

func Run() error {

	redisClient, err := redis.NewPubSubClient(nil)
	if err != nil {
		return err
	}

	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		clusterName = "test_cluster"
	}

	ch := redisClient.Subscribe(context.Background(), clusterName)

	l, err := listener.NewListener(ch, clusterName)
	if err != nil {
		return err
	}

	l.ListenAndServe()
	return nil // should never reach here
}
