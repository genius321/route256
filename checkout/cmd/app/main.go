package main

import (
	"log"
	"net/http"
	"route256/checkout/clients/loms"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	model := domain.New(loms.New("http://localhost:8081"))

	hand := &addtocart.Handler{
		Model: model,
	}

	http.Handle("/addToCart", srvwrapper.New(hand.Handle))

	err := http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
