package response

type UserResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}
