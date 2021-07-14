package rqctx

import (
	"CloudScapes/pkg/wire"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Context) DecodeBody(out interface{}) error {
	if err := json.NewDecoder(c.r.Body).Decode(&out); err != nil {
		return wire.NewBadRequestError(err.Error())
	}
	return nil
}

func (c *Context) IdFromPath(key string) (int64, error) {
	vars := mux.Vars(c.r)
	val, ok := vars[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' was not found at path", key)
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}
