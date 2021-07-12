package apihandlers

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/wire"
)

func UsersGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Users.GetAllUsers()
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
