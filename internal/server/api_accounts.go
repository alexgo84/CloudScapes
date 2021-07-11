package server

import "CloudScapes/internal/server/rqctx"

func accountsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	return c.SendNothing()
}

func accountsPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	return c.SendNothing()
}
