package request

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email,uniqueEmail"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8"`
	Phone                string `json:"phone"`
	Address              string `json:"address"`
}

type UserUpdateRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email,uniqueEmail"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdatePasswordRequest struct {
	UserID          string `json:"user_id"`
	NewPassword     string `json:"password_new" validate:"required,min=8"`
	ConfirmPassword string `json:"password_confirmation" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=8,required"`
}
