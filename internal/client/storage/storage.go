package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	// read file
	ubytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// parse file data
	if len(ubytes) != 0 {
		return json.Unmarshal(ubytes, &Users)
	}

	return nil
}

// initLocal reads or creates the local users data storage file. Then parse the content to local memory.
func initLocal(filePath string) error {
	// second - open vault file
	//Objects := make(map[string]*model.LocalStorage)

	// create/open file
	fo, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fo.Close()

	// read the whole file at once
	vbytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// parse file data
	if len(vbytes) != 0 {
		return json.Unmarshal(vbytes, &Objects)
	}

	return nil
}

// UpdateFiles rewrites local files with actual data.
func UpdateFiles(userFilePath, objectFilePath string) error {
	// prepare users data
	usersJSONByte, err := json.Marshal(Users)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = UpdateFile(userFilePath, usersJSONByte); err != nil {
		return err
	}

	// prepare vault data
	vaultJSONByte, err := json.Marshal(Objects)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = UpdateFile(objectFilePath, vaultJSONByte); err != nil {
		return err
	}

	return nil
}

// UpdateFile method rewrite the file with the passed data.
func UpdateFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
