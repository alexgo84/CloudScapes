package wire

type APIError struct {
	StatusCode int
	Err        error
}

func (e APIError) Error() string {
	return e.Err.Error()
}
