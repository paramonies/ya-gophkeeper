package model

import (
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

func ProtoToLocalStorage(in *pb.GetAllUserDataFromDBResponse) *LocalStorage {
	return &LocalStorage{
		Password: createPasswordMap(in.Passwords),
		Text:     createTextMap(in.Texts),
		Binary:   createBinaryMap(in.Binaries),
		Card:     createCardMap(in.Cards),
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

func createTextMap(texts []*pb.Text) map[string]*Text {
	out := make(map[string]*Text)

	for _, t := range texts {
		out[t.GetTitle()] = &Text{
			Title:   t.GetTitle(),
			Data:    t.GetData(),
			Meta:    t.GetMeta(),
			Version: t.GetVersion(),
		}
	}

	return out
}

func createBinaryMap(bins []*pb.Binary) map[string]*Binary {
	out := make(map[string]*Binary)

	for _, b := range bins {
		out[b.GetTitle()] = &Binary{
			Title:   b.GetTitle(),
			Data:    b.GetData(),
			Meta:    b.GetMeta(),
			Version: b.GetVersion(),
		}
	}

	return out
}

func createCardMap(cards []*pb.Card) map[string]*Card {
	out := make(map[string]*Card)

	for _, c := range cards {
		out[c.GetNumber()] = &Card{
			Number:  c.GetNumber(),
			Owner:   c.GetOwner(),
			ExpDate: c.GetExpDate(),
			Cvv:     c.GetCvv(),
			Meta:    c.GetMeta(),
			Version: c.GetVersion(),
		}
	}

	return out
}
