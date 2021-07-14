package wire

type NewUser struct {
	Email    string `json:"email" db:"email"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password,omitempty" db:"-"`
}
