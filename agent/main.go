package main

import (
	"CloudScapes/pkg/logger"
	"errors"
)

func main() {
	if err := logger.InitLogger(true, nil); err != nil {
		panic(err)
	}
	defer logger.Flush()

	panic(errors.New("agent not implemented yet"))
}
