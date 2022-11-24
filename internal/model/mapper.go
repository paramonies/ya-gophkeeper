package model

import (
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

func PasswordModelstoProto(pwds []*Password) []*pb.Password {
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
