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
}

type Binary struct {
}

type Card struct {
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
