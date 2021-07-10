package rqctx

import (
	"context"

	"github.com/google/uuid"
)

type Context struct {
	ctx  context.Context
	uuid uuid.UUID
}

func NewRequestContext(ctx context.Context) Context {
	return Context{
		ctx:  ctx,
		uuid: uuid.New(),
	}
}
