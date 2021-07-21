package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const PubSubChannelSize = 1000

type Credentials struct {
	hostname string
	password string
	port     int
	tls      bool
}

func NewCredentials(host, pass string, port int, tls bool) Credentials {
	return Credentials{
		hostname: host,
		port:     port,
		password: pass,
		tls:      tls,
	}
}

func (c Credentials) Address() string {
	return fmt.Sprintf("%s:%d", c.hostname, c.port)
}

func (c Credentials) Password() string {
	return c.password
}

func (c *Credentials) getTLSConfig() *tls.Config {
	if c.tls {
		return &tls.Config{}
	}
	return nil
}

type PubSubClient struct {
	redisClient *redis.Client
}

func NewPubSubClient(creds Credentials) (*PubSubClient, error) {

	// Create a new Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:      creds.Address(),
		Password:  creds.Password(),
		DB:        0,
		TLSConfig: creds.getTLSConfig(),
	})

	// Ping the Redis server and check if any errors occured
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			return nil, err
		}
	}

	return &PubSubClient{
		redisClient: redisClient,
	}, nil
}

func (ps *PubSubClient) Subscribe(ctx context.Context, channel string) <-chan *redis.Message {
	ch := ps.redisClient.Subscribe(ctx, channel)
	ch.ChannelSize(PubSubChannelSize)
	return ch.Channel()
}

func (ps *PubSubClient) Publish(ctx context.Context, channel string, message []byte) error {
	rv := ps.redisClient.Publish(ctx, channel, message)
	return rv.Err()
}

func (ps *PubSubClient) Close() error {
	return ps.redisClient.Close()
}
