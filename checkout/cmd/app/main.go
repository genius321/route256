package main

import (
	"log"
	"net/http"
	"route256/checkout/cmd/internal/domain"
	"route256/checkout/cmd/internal/handlers/addtocart"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	model := &domain.Model{}

	hand := &addtocart.Handler{
		Model: model,
	}

	http.Handle("/addToCart", srvwrapper.New(hand.Handle))

	err := http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
