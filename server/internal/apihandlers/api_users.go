package apihandlers

import (
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/rqctx"
)

func UsersGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Users.GetUsers()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(accounts)
}

func UsersPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	var newUser wire.NewUser
	if err := c.DecodeBody(&newUser); err != nil {
		return c.SendError(err)
	}

	user, err := c.Users.CreateUser(&newUser)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}
