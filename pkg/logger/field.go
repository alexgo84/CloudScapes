package logger

import "go.uber.org/zap"

// Str is acting as a facade for a zap field. use in order to not import zap where not needed
func Str(k string, v string) zap.Field {
	return zap.String(k, v)
}

// Int64 is acting as a facade for a zap field. use in order to not import zap where not needed
func Int64(k string, v int64) zap.Field {
	return zap.Int64(k, v)
}

// Err is acting as a facade for a zap field. use in order to not import zap where not needed
func Err(err error) zap.Field {
	return zap.Error(err)
}
