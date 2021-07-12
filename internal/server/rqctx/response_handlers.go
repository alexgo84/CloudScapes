package rqctx

import (
	"CloudScapes/pkg/wire"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type ResponseHandler struct {
	StatusCode int
	Obj        interface{}
	Err        error
}

func (c *Context) SendOK(obj interface{}) ResponseHandler {
	return ResponseHandler{
		StatusCode: http.StatusOK,
		Obj:        obj,
		Err:        nil,
	}
}

func (c *Context) SendCreated(obj interface{}) ResponseHandler {
	return ResponseHandler{
		StatusCode: http.StatusCreated,
		Obj:        obj,
		Err:        nil,
	}
}

func (c *Context) SendNothing() ResponseHandler {
	return ResponseHandler{
		StatusCode: http.StatusNoContent,
		Obj:        nil,
		Err:        nil,
	}
}

func (c *Context) SendError(err error) ResponseHandler {
	var rv ResponseHandler
	var apiErr wire.APIError
	if ok := errors.As(err, &apiErr); ok {
		rv = ResponseHandler{
			StatusCode: apiErr.StatusCode,
			Obj:        nil,
			Err:        apiErr.Err,
		}
	} else {
		rv = ResponseHandler{
			StatusCode: http.StatusInternalServerError,
			Obj:        nil,
			Err:        err,
		}
	}

	if rv.StatusCode >= 500 {
		c.ReportError("internal server error", zap.Error(err))
	}
	return rv
}
