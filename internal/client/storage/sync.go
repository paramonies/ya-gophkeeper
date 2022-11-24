package storage

import "github.com/paramonies/ya-gophkeeper/internal/model"

// SyncData check for latest version from local storage and DB. Returns storage with the latest versions from both.
func SyncData(dataLocal, dataDB *model.LocalStorage) *model.LocalStorage {
	out := CreateStorage()

	//anti nil incoming data.
	if dataLocal == nil {
		dataLocal = CreateStorage()
	}
	if dataDB == nil {
		dataDB = CreateStorage()
	}

	// save key-titles from both structures
	entitiesP := getAllEntities("password", dataLocal, dataDB)
	entitiesT := getAllEntities("text", dataLocal, dataDB)
	entitiesB := getAllEntities("binary", dataLocal, dataDB)
	entitiesC := getAllEntities("card", dataLocal, dataDB)

	for e := range entitiesP {
		out.Password[e] = FindLatestPassword(e, dataLocal.Password, dataDB.Password)
	}

	for e := range entitiesT {
		out.Text[e] = FindLatestText(e, dataLocal.Text, dataDB.Text)
	}

	for e := range entitiesB {
		out.Binary[e] = FindLatestBinary(e, dataLocal.Binary, dataDB.Binary)
	}

	for e := range entitiesC {
		out.Card[e] = FindLatestCard(e, dataLocal.Card, dataDB.Card)
	}

	return out
}

func getAllEntities(dataType string, s1, s2 *model.LocalStorage) map[string]struct{} {
	entities := map[string]struct{}{}

	switch dataType {
	case "password":
		for k := range s1.Password {
			entities[k] = struct{}{}
		}
		for k := range s2.Password {
			entities[k] = struct{}{}
		}
	case "text":
		for k := range s1.Text {
			entities[k] = struct{}{}
		}
		for k := range s2.Text {
			entities[k] = struct{}{}
		}
	case "bin":
		for k := range s1.Binary {
			entities[k] = struct{}{}
		}
		for k := range s2.Binary {
			entities[k] = struct{}{}
		}
	case "card":
		for k := range s1.Card {
			entities[k] = struct{}{}
		}
		for k := range s2.Card {
			entities[k] = struct{}{}
		}
	}

	return entities
}

func FindLatestPassword(e string, dataLocal, dataDB map[string]*model.Password) *model.Password {
	l, lOK := dataLocal[e]
	db, dbOK := dataDB[e]
	if lOK && dbOK {
		if l.Version > db.Version {
			return l
		} else {
			return db
		}
	}

	if lOK && !dbOK {
		return l
	}
	if !lOK && dbOK {
		return db
	}

	return nil
}

func FindLatestText(e string, dataLocal, dataDB map[string]*model.Text) *model.Text {
	l, lOK := dataLocal[e]
	db, dbOK := dataDB[e]
	if lOK && dbOK {
		if l.Version > db.Version {
			return l
		} else {
			return db
		}
	}

	if lOK && !dbOK {
		return l
	}
	if !lOK && dbOK {
		return db
	}

	return nil
}

func FindLatestBinary(e string, dataLocal, dataDB map[string]*model.Binary) *model.Binary {
	l, lOK := dataLocal[e]
	db, dbOK := dataDB[e]
	if lOK && dbOK {
		if l.Version > db.Version {
			return l
		} else {
			return db
		}
	}

	if lOK && !dbOK {
		return l
	}
	if !lOK && dbOK {
		return db
	}

	return nil
}

func FindLatestCard(e string, dataLocal, dataDB map[string]*model.Card) *model.Card {
	l, lOK := dataLocal[e]
	db, dbOK := dataDB[e]
	if lOK && dbOK {
		if l.Version > db.Version {
			return l
		} else {
			return db
		}
	}

	if lOK && !dbOK {
		return l
	}
	if !lOK && dbOK {
		return db
	}

	return nil
}
