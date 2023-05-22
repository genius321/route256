package purchase

import (
	"errors"
	"log"
)

type Handler struct {
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Request struct {
	User int64 `json:"user"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (h Handler) Handle(req Request) (Response, error) {
	log.Printf("%+v", req)
	return Response{}, nil
}
