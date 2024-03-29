package rqctx

import (
	"CloudScapes/pkg/wire"
	"errors"
	"net/http"
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
			Err:        apiErr,
		}
	} else {
		rv = ResponseHandler{
			StatusCode: http.StatusInternalServerError,
			Obj:        nil,
			Err:        err,
		}
	}

	return rv
}
