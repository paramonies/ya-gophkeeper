package storage

import (
	"encoding/json"
	"os"

	"github.com/paramonies/ya-gophkeeper/internal/model"
)

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

// InitStorage function initializes the storage data (check files & parse to local memory).
func InitStorage(userFilePath, objectFilePath string) (err error) {
	if err = initUsers(userFilePath); err != nil {
		return err
	}

	if err = initLocal(objectFilePath); err != nil {
		return err
	}

	return nil
}

// initUsers reads or creates the local user auth info file. Then parse the content to local memory.
func initUsers(filePath string) error {
	// create/open file
	fu, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer fu.Close()

	ubytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(ubytes) != 0 {
		return json.Unmarshal(ubytes, &Users)
	}

	return nil
}

// initLocal reads or creates the local users data storage file. Then parse the content to local memory.
func initLocal(filePath string) error {
	// create/open file
	fo, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer fo.Close()

	vbytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(vbytes) != 0 {
		return json.Unmarshal(vbytes, &Objects)
	}

	return nil
}

// UpdateFiles rewrites local files with actual data.
func UpdateFiles(userFilePath, objectFilePath string) error {
	usersJSONByte, err := json.Marshal(Users)
	if err != nil {
		return err
	}
	if err = UpdateFile(userFilePath, usersJSONByte); err != nil {
		return err
	}

	objectsJSONByte, err := json.Marshal(Objects)
	if err != nil {
		return err
	}

	if err = UpdateFile(objectFilePath, objectsJSONByte); err != nil {
		return err
	}

	return nil
}

// UpdateFile method rewrite the file with the passed data.
func UpdateFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
