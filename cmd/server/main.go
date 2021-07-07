package main

import (
	"CloudScapes/pkg/logger"

	"CloudScapes/internal/server"

	"go.uber.org/zap"
)

func main() {

	if err := logger.InitLogger(true, nil); err != nil {
		panic(err)
	}
	// since we log both to console and file we must sync. we ignore error
	// as it will always complain since console is not syncable
	// https://github.com/uber-go/zap/issues/880
	defer zap.L().Sync()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
