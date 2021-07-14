package apihandlers

import (
	"CloudScapes/pkg/logger"
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/dat"
	"CloudScapes/server/internal/rqctx"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

func AccountsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	accounts, err := c.Accounts.GetAccounts()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(accounts)
}

// AccountsPostHandler will create an account if one does not exist. Since there is no account here
// and the endpoint is not authenticated, we handle request/response flow differently
func AccountsPostHandler(db *dat.DB) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		shouldCommit := false
		sendAndReportInternalError := func(err error) {
			respond(rw, http.StatusInternalServerError, "Internal server error")
			logger.Log(logger.ERROR, "failed to create account", logger.Err(err))
		}

		ctx := context.TODO()
		txn, err := db.GetNewTransaction(ctx)
		if err != nil {
			sendAndReportInternalError(err)
			return
		}

		defer func() {
			if shouldCommit {
				if err := txn.Commit(); err != nil {
					logger.Log(logger.ERROR, "failed to commit account creation", logger.Err(err))
				}
			} else {
				if err := txn.Rollback(); err != nil {
					logger.Log(logger.ERROR, "failed to rollback account creation", logger.Err(err))
				}
			}
		}()

		aMapper := dat.NewAccountsMapper(ctx, txn)
		accounts, err := aMapper.GetAccounts()
		if err != nil {
			sendAndReportInternalError(err)
			return
		}

		if IsProduction() && len(accounts) >= 1 {
			respond(rw, http.StatusBadRequest, "only one account may be created in a production environment")
			return
		}

		var req wire.CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendAndReportInternalError(err)
			return
		}

		a, err := aMapper.CreateAccount(req.CompanyName)
		if err != nil {
			sendAndReportInternalError(err)
			return
		}

		newUser := wire.NewUser{
			Name:     req.Email,
			Email:    req.Email,
			Password: req.Password,
		}
		uMapper := dat.NewUsersMapper(ctx, txn, a.ID)
		u, err := uMapper.CreateUser(&newUser)
		if err != nil {
			sendAndReportInternalError(err)
			return
		}

		// if responding with 201 was successful we should commit the transaction
		shouldCommit = respond(rw, http.StatusCreated, u)
	}
}

func IsProduction() bool {
	return os.Getenv("RELEASE_STAGE") == "production"
}
