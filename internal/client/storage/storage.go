package storage

import "github.com/paramonies/ya-gophkeeper/internal/model"

var (
	Users   = make(map[string]string)
	Objects = make(map[string]*model.LocalStorage)
)

func CreateStorage() *model.LocalStorage {
	return &model.LocalStorage{
		Password: make(map[string]*model.Password),
		Text:     make(map[string]*model.Text),
		Binary:   make(map[string]*model.Binary),
		Card:     make(map[string]*model.Card),
	}
}
