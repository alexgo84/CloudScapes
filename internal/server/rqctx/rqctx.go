package rqctx

import (
	"CloudScapes/internal/server/dat"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler func(c *Context) ResponseHandler

type Context struct {
	r      *http.Request
	writer http.ResponseWriter
	txn    *sqlx.Tx
	uuid   uuid.UUID
	dat.DataContext
}

func NewRequestContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		r:      r,
		uuid:   uuid.New(),
		writer: w,
	}
}

func (ctx *Context) InitDBContext(db *dat.DB) error {
	txn, err := db.GetNewTransaction(ctx.r.Context())
	if err != nil {
		return err
	}
	ctx.txn = txn
	ctx.DataContext = dat.NewDataContext(ctx.r.Context(), txn)
	return nil
}

func (ctx *Context) Commit() error {
	return ctx.txn.Commit()
}

func (ctx *Context) Rollback() error {
	return ctx.txn.Rollback()
}
