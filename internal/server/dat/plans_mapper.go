package dat

import (
	"CloudScapes/pkg/wire"
	"time"
)

type Plan struct {
	ID      int64     `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created_at"`
	wire.Plan
}
