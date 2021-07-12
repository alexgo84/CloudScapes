package rqctx

import (
	"CloudScapes/pkg/wire"
	"encoding/json"
	"net/http"
)

func (c *Context) DecodeBody(out interface{}) error {
	if err := json.NewDecoder(c.r.Body).Decode(&out); err != nil {
		apiErr := wire.APIError{Err: err, StatusCode: http.StatusBadRequest}
		return apiErr
	}
	return nil
}
