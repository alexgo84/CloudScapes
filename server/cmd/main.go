package main

import (
	"CloudScapes/pkg/logger"

	server "CloudScapes/server/internal"
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
