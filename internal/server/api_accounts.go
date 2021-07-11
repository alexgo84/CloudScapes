package server

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/wire"
	"errors"
)

func accountsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Accounts.GetAccounts()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(accounts)
}

func accountsPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Accounts.GetAccounts()
	if err != nil {
		return c.SendError(err)
	}
	if len(accounts) >= 1 {
		multipleAccountsError := wire.APIError{
			StatusCode: 400,
			Err:        errors.New("only one account may be created"),
		}
		c.SendError(multipleAccountsError)
	}

	account, err := c.Accounts.CreateAccount()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(account)
}
