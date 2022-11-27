package service

import (
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var (
	newerVersionDetected = "newer version found in database. please synchronize you app to get the most actual data."
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

func TextModelsToProto(texts []*model.Text) []*pb.Text {
	res := make([]*pb.Text, len(texts))

	for _, t := range texts {
		r := &pb.Text{
			Title:   t.Title,
			Data:    t.Data,
			Meta:    t.Meta,
			Version: t.Version,
		}

		res = append(res, r)
	}
	return res
}

func BinaryModelsToProto(binaries []*model.Binary) []*pb.Binary {
	res := make([]*pb.Binary, len(binaries))

	for _, b := range binaries {
		b := &pb.Binary{
			Title:   b.Title,
			Data:    b.Data,
			Meta:    b.Meta,
			Version: b.Version,
		}

		res = append(res, b)
	}
	return res
}

func CardModelsToProto(cards []*model.Card) []*pb.Card {
	res := make([]*pb.Card, len(cards))

	for _, card := range cards {
		c := &pb.Card{
			Number:  card.Number,
			Owner:   card.Owner,
			ExpDate: card.ExpDate,
			Cvv:     card.Cvv,
			Meta:    card.Meta,
			Version: card.Version,
		}

		res = append(res, c)
	}
	return res
}
