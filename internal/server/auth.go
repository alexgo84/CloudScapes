package server

import (
	"CloudScapes/internal/server/rqctx"
)

func authSession(h rqctx.Handler) rqctx.Handler {
	return func(c *rqctx.Context) rqctx.ResponseHandler {
		return h(c)
	}
}
