package main

import (
	"CloudScapes/pkg/logger"

	"CloudScapes/internal/server"
)

func main() {

	if err := logger.InitLogger(true, nil); err != nil {
		panic(err)
	}
	defer logger.Flush()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
