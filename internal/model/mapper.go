package model

import (
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

func ProtoToLocalStorage(in *pb.GetAllUserDataFromDBResponse) *LocalStorage {
	return &LocalStorage{
		Password: createPasswordMap(in.Passwords),
		// todo: make for other entities
	}
}

func createPasswordMap(pwds []*pb.Password) map[string]*Password {
	out := make(map[string]*Password)

	for _, p := range pwds {
		out[p.GetLogin()] = &Password{
			Login:    p.GetLogin(),
			Password: p.GetPassword(),
			Meta:     p.GetMeta(),
			Version:  p.GetVersion(),
		}
	}

	return out
}

// todo: createMap make for other entities
