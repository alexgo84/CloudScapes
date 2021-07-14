package wire

type NewCluster struct {
	Name string `json:"name" db:"name"`
}
