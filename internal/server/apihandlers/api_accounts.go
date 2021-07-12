package apihandlers

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/wire"
	"errors"
	"fmt"
)

func AccountsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Accounts.GetAccounts()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(accounts)
}

func AccountsPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Accounts.GetAccounts()
	if err != nil {
		return c.SendError(err)
	}

	if len(accounts) >= 1 {
		multipleAccountsError := wire.APIError{
			StatusCode: 400,
			Err:        errors.New("only one account may be created"),
		}
		return c.SendError(multipleAccountsError)
	}

	var req wire.PostAccountRequest
	if err := c.DecodeBody(&req); err != nil {
		return c.SendError(err)
	}

	account, err := c.Accounts.CreateAccount(req.CompanyName)
	if err != nil {
		return c.SendError(err)
	}

	fmt.Println(account)

	return c.SendCreated(account)
}
