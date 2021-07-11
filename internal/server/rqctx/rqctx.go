package rqctx

import (
	"CloudScapes/internal/server/dat"
	"CloudScapes/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Handler func(c *Context) ResponseHandler

type Context struct {
	r      *http.Request
	writer http.ResponseWriter
	txn    pgx.Tx
	uuid   uuid.UUID
	dat.DataContext
}

func NewRequestContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		r:      r,
		uuid:   uuid.New(),
		writer: w,
	}
}

func (ctx *Context) InitDBTransaction() error {
	txn, err := dat.GetNewTransaction(ctx.r.Context())
	if err != nil {
		return err
	}
	ctx.txn = txn
	ctx.DataContext = dat.NewDataContext(txn)
	return nil
}

type Field zap.Field

// ReportError will send an error to a tracking tool after enriching it with the contexts data
func (ctx *Context) ReportError(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("uuid", ctx.uuid.String()))
	zap.L().Error(msg, fields...)
}

func (ctx *Context) SendAPIError(err error) {
	ctx.ReportError(err.Error())

	ctx.writer.Header().Set("Content-Type", "application/json")
	marshalledError, err := json.Marshal(err)
	if err != nil {
		http.Error(ctx.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = ctx.writer.Write(marshalledError); err != nil {
		ctx.ReportError(err.Error())
		http.Error(ctx.writer, err.Error(), http.StatusInternalServerError)
	}

}

func (c *Context) MarshalAndWrite(payload interface{}, status int) []byte {
	if payload == nil || status == http.StatusNoContent {
		c.writer.WriteHeader(status)
		return []byte{}
	}

	marshaled, err := json.Marshal(payload)
	if err != nil {
		c.writer.WriteHeader(http.StatusInternalServerError)
		response := []byte("{\"message\":\"Failed to serialize response\"}")
		logger.Log(logger.ERROR, "failed to serialize response", logger.Err(err))

		if _, err := c.writer.Write(response); err != nil {
			logger.Log(logger.ERROR, "failed to write serialization failure response", logger.Err(err))
		}

		return response
	}

	c.writer.Header().Set("Content-Type", "application/json")
	c.writer.Header().Set("X-Frame-Options", "SAMEORIGIN")

	c.writer.WriteHeader(status)

	if _, err := c.writer.Write(marshaled); err != nil {
		logger.Log(logger.ERROR, "failed to write serialized response", logger.Err(err))
	}
	return marshaled
}

func (ctx *Context) Commit() error {
	return ctx.txn.Commit(ctx.r.Context())
}

func (ctx *Context) Rollback() error {
	return ctx.txn.Rollback(ctx.r.Context())
}
