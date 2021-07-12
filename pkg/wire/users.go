package wire

type NewUser struct {
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	AccountID int64  `json:"accountId" db:"accountid"`
	Password  string `json:"password,omitempty" db:"-"`
}
