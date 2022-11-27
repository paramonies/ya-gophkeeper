package model

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Password struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
	Version  uint32 `json:"version"`
}

type Text struct {
	Title   string `json:"title"`
	Data    string `json:"data"`
	Meta    string `json:"meta"`
	Version uint32 `json:"version"`
}

type Binary struct {
	Title   string `json:"title"`
	Data    string `json:"data"`
	Meta    string `json:"meta"`
	Version uint32 `json:"version"`
}

type Card struct {
	Number  string `json:"number"`
	Owner   string `json:"owner"`
	ExpDate string `json:"expiration_date"`
	Cvv     string `json:"cvv"`
	Meta    string `json:"meta"`
	Version uint32 `json:"version"`
}

// LocalStorage is a local struct for client data
type LocalStorage struct {
	Password map[string]*Password `json:"password"`
	Text     map[string]*Text     `json:"text"`
	Binary   map[string]*Binary   `json:"binary"`
	Card     map[string]*Card     `json:"card"`
}

// Server data
type ServerData struct {
	Passwords []*Password `json:"passwords"`
	Texts     []*Text     `json:"texts"`
	Binaries  []*Binary   `json:"binaries"`
	Cards     []*Card     `json:"cards"`
}
