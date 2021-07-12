package wire

type CreateAccountRequest struct {
	CompanyName string `json:"companyName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}
