package wire

import (
	"net/http"
)

type APIError struct {
	StatusCode      int
	Message         string
	UnderlyingError error
}

func (e APIError) Error() string {
	return e.UnderlyingError.Error()
}

func NewBadRequestError(message string) APIError {
	return APIError{StatusCode: http.StatusBadRequest, Message: message}
}

func NewConflictError(message string, err error) APIError {
	return APIError{StatusCode: http.StatusConflict, Message: message, UnderlyingError: err}
}

func NewNotFoundError(message string, err error) APIError {
	return APIError{StatusCode: http.StatusNotFound, Message: message, UnderlyingError: err}
}
