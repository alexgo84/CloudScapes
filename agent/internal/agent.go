package agent

import (
	"CloudScapes/agent/internal/listener"
	"CloudScapes/pkg/shared/redis"
	"context"
	"os"
)

func Run() error {

	creds := redis.NewCredentials("localhost", "", 6379, false)
	redisClient, err := redis.NewPubSubClient(creds)
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

	return l.ListenAndServe()
}
