package dat

import (
	"CloudScapes/pkg/wire"
	"time"
)

type Deployment struct {
	ID        int64     `json:"id" db:"id"`
	Created   time.Time `json:"created" db:"created_at"`
	CreatedBy int64     `json:"createdBy" db:"created_by"`

	Modified   *time.Time `json:"modified" db:"modfied_at"`
	ModifiedBy *int64     `json:"modifiedBy" db:"modified_by"`
	Deleted    *time.Time `json:"deleted" db:"deleted_at"`
	DeletedBy  *int64     `json:"deletedBy" db:"deleted_by"`

	SalesforceState *string `json:"salesforceState" db:"sf_state"`

	wire.Deployment
}
