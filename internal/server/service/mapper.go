package service

import (
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

func PasswordModelsToProto(pwds []*model.Password) []*pb.Password {
	res := make([]*pb.Password, len(pwds))

	for _, pwd := range pwds {
		r := &pb.Password{
			Login:    pwd.Login,
			Password: pwd.Password,
			Meta:     pwd.Meta,
			Version:  pwd.Version,
		}

		res = append(res, r)
	}
	return res
}

// todo: make modelsToProto for other entities
