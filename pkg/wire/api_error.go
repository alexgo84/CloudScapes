package wire

import (
	"errors"
	"net/http"
)

type APIError struct {
	StatusCode int
	Err        error
}

func (e APIError) Error() string {
	return e.Err.Error()
}

func NewBadRequestError(message string) APIError {
	err := errors.New(message)
	return APIError{StatusCode: http.StatusBadRequest, Err: err}
}

func NewConflictError(message string) APIError {
	err := errors.New(message)
	return APIError{StatusCode: http.StatusConflict, Err: err}
}
