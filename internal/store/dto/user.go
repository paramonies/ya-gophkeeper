package dto

type RegisterRequest struct {
	Login        string
	PasswordHash string
}

type RegisterResponse struct {
	UserID string
}
