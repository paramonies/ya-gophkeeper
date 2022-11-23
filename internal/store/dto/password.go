package dto

type CreateRequest struct {
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type CreateResponse struct {
	PasswordID string
}

type GetByIDRequest struct {
	Login  string
	UserID string
}

type GetByIDResponse struct {
	ID       string
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type GetPasswordRequest struct {
}

type GetPasswordResponse struct {
}

type UpdatePasswordRequest struct {
}

type UpdatePasswordResponse struct {
}

type DeletePasswordRequest struct {
}

type DeletePasswordResponse struct {
}
