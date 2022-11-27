package dto

type RegisterRequest struct {
	Login        string
	PasswordHash string
}

type RegisterResponse struct {
	UserID string
}

type LoginRequest struct {
	Login string
}

type LoginResponse struct {
	UserID       string
	PasswordHash string
}
