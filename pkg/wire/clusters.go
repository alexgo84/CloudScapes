package wire

type NewCluster struct {
	AccountID int64 `json:"accountId" db:"accountid"`

	Name string `json:"name" db:"name"`
}
