package server

import (
	"CloudScapes/pkg/logger"
	l "CloudScapes/pkg/logger"
	"CloudScapes/pkg/pubsub"
	"CloudScapes/pkg/shared"
	"CloudScapes/server/internal/dat"
	"CloudScapes/server/internal/rqctx"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

func Run() error {
	l.Log(l.INFO, "initializing database")

	pubsub, err := pubsub.NewPubSubClient(nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := pubsub.Close(); err != nil {
			l.Log(l.ERROR, "pubsub close failed", l.Err(err))
		}
	}()

	db, err := dat.NewDB(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Log(l.ERROR, "database close failed", l.Err(err))
		}
	}()

	l.Log(l.INFO, "running all database migrations")
	if err := db.RunMigrations(context.Background()); err != nil {
		return err
	}

	go startListeningForControlMessages(pubsub, db)

	s := &http.Server{
		Handler:      createRouter(db, pubsub),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	l.Log(l.INFO, "serving requests on port 8080")
	return s.ListenAndServe()
}

func handlePanicRecovery(ctx *rqctx.Context) {
	if err := recover(); err != nil {

		// report error
		e, ok := err.(error)
		if !ok {
			e = fmt.Errorf("%#v\n%s", err, string(debug.Stack()))
		}
		ctx.ReportError("Recovered handler panic", zap.Error(e))

		// print stack
		l.Log(l.INFO, string(debug.Stack()))

		// rollback transaction
		rErr := ctx.Rollback()
		if err != nil {
			l.Log(l.ERROR, "failed to rollback error", zap.Error(rErr))
		}
	}
}

func startListeningForControlMessages(pubsub *pubsub.PubSubClient, db *dat.DB) {
	ctx := context.Background()
	ch := pubsub.Subscribe(ctx, shared.ControlChannelName)
	for msg := range ch {
		payload := []byte(msg.Payload)
		var agentRes shared.AgentResponse
		if err := json.Unmarshal(payload, &agentRes); err != nil {
			logger.Log(logger.ERROR, "failed to read control message", logger.Err(err))
			continue
		}

		if err := handleAgentResponse(&agentRes); err != nil {
			logger.Log(logger.ERROR, "failed to handle control message", logger.Err(err))
		}
	}
}

func handleAgentResponse(res *shared.AgentResponse) error {
	logger.Log(logger.INFO, "agent says", logger.Any("response", res))
	return nil
}
