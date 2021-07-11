package server

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/logger"
)

func accountsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	logger.Log(logger.INFO, "IN accountsGetHandler")
	return c.SendNothing()
}

func accountsPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	logger.Log(logger.INFO, "IN accountsPostHandler")
	return c.SendNothing()
}
