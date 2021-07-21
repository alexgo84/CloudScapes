package main

import (
	agent "CloudScapes/agent/internal"
	"CloudScapes/pkg/logger"
)

func main() {
	if err := logger.InitLogger(true, nil); err != nil {
		panic(err)
	}
	defer logger.Flush()

	if err := agent.Run(); err != nil {
		panic(err)
	}
}
