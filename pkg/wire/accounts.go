package wire

type PostAccountRequest struct {
	CompanyName string `json:"companyName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}
