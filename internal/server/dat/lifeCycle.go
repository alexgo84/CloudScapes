package dat

import (
	"CloudScapes/pkg/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

var dbMap *pgx.Conn

func InitDB(ctx context.Context) error {
	user := getEnv("POSTGRES_USER", "cloudscapes")
	pass := getEnv("POSTGRES_DB", "cloudscapes")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	dbName := getEnv("POSTGRES_DB_NAME", "cloudscapes")

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbName)

	var err error
	dbMap, err = pgx.Connect(ctx, connectionURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	logger.Log(logger.INFO, fmt.Sprintf("Connected user '%s' to DB '%s' @ %s:%s", user, dbName, host, port))

	return nil
}

func CloseDB(ctx context.Context) error {
	return dbMap.Close(ctx)
}

func PingDB(ctx context.Context) error {
	if dbMap == nil {
		return errors.New("db not initialized")
	}
	contextWithTimeout, c := context.WithTimeout(ctx, time.Second*5)
	err := dbMap.Ping(contextWithTimeout)
	c()
	return err
}

func GetNewTransaction(ctx context.Context) (pgx.Tx, error) {
	return dbMap.Begin(ctx)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
