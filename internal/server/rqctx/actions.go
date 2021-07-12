package rqctx

import (
	"CloudScapes/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ReportError will send an error to a tracking tool after enriching it with the contexts data
func (ctx *Context) ReportError(msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("uuid", ctx.uuid.String()))
	logger.Log(logger.ERROR, msg, fields...)
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
