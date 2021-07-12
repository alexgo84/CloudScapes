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

		defer func() {
			handlePanicRecovery(ctx)
		}()

		if err := ctx.InitDBContext(db); err != nil {
			logger.Log(logger.INFO, "failed to init DB context", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()))
			if err := ctx.Rollback(); err != nil {
				ctx.ReportError("rollback failed", logger.Err(err))
			}
			ctx.MarshalAndWrite([]byte("Internal server error"), http.StatusInternalServerError)
			return
		}

		res := h(ctx)

		switch {
		case res.StatusCode >= 400:
			logger.Log(logger.INFO, "write ERROR response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Rollback(); err != nil {
				ctx.ReportError("rollback failed", zap.Error(err))
			}
			ctx.MarshalAndWrite([]byte(res.Err.Error()), res.StatusCode)
			return

		default:
			logger.Log(logger.INFO, "write response", logger.Str("method", r.Method), logger.Str("URL", r.URL.String()), logger.Int64("status", int64(res.StatusCode)))
			if err := ctx.Commit(); err != nil {
				ctx.ReportError("commit failed", zap.Error(err))
				return
			}
			ctx.MarshalAndWrite(res.Obj, res.StatusCode)
			return
		}
	}
}

func authSession(h rqctx.Handler) rqctx.Handler {
	return func(c *rqctx.Context) rqctx.ResponseHandler {
		return h(c)
	}
}
