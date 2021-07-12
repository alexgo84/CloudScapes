package server

import (
	"CloudScapes/internal/server/dat"
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

func contextify(db *dat.DB, h rqctx.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := rqctx.NewRequestContext(w, r)

		ctx.InitDBContext(db)

		defer func() {
			handlePanicRecovery(ctx)
		}()

		res := h(ctx)

		switch {
		case res.StatusCode >= 400:
			logger.Log(logger.INFO, "write ERROR response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Rollback(); err != nil {
				logger.Log(logger.ERROR, "commit failed", zap.Error(err))
			}
			ctx.MarshalAndWrite([]byte(res.Err.Error()), res.StatusCode)

		default:
			logger.Log(logger.INFO, "write response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Commit(); err != nil {
				logger.Log(logger.ERROR, "commit failed", zap.Error(err))
				return
			}
			ctx.MarshalAndWrite(res.Obj, res.StatusCode)
		}
	}
}

func authSession(h rqctx.Handler) rqctx.Handler {
	return func(c *rqctx.Context) rqctx.ResponseHandler {
		return h(c)
	}
}
