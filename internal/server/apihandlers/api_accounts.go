package apihandlers

import (
	"CloudScapes/internal/server/rqctx"
	"CloudScapes/pkg/wire"
	"errors"
	"os"
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

	if IsProduction() && len(accounts) >= 1 {
		multipleAccountsError := wire.APIError{
			StatusCode: 400,
			Err:        errors.New("only one account may be created in a production environment"),
		}
		return c.SendError(multipleAccountsError)
	}

	var req wire.CreateAccountRequest
	if err := c.DecodeBody(&req); err != nil {
		return c.SendError(err)
	}

	account, err := c.Accounts.CreateAccount(req.CompanyName)
	if err != nil {
		return c.SendError(err)
	}

	newUser := wire.NewUser{
		Name:      req.Email,
		Email:     req.Email,
		AccountID: account.ID,
		Password:  req.Password,
	}
	if _, err := c.Users.CreateUser(&newUser); err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(newUser)
}

func IsProduction() bool {
	return os.Getenv("RELEASE_STAGE") == "production"
}
