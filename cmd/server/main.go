package main

import (
	"CloudScapes/pkg/logger"

	"go.uber.org/zap"
)

func main() {

	if err := logger.InitLogger(true); err != nil {
		panic(err)
	}
	// since we log both to console and file we must sync. we ignore error
	// as it will always complain since console is not syncable
	// https://github.com/uber-go/zap/issues/880
	defer zap.L().Sync()
}
