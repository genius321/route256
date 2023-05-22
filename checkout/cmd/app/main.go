package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	err := config.Init()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	model := domain.New(loms.New(config.AppConfig.Services.Loms))

	hand := &addtocart.Handler{
		Model: model,
	}

	http.Handle("/addToCart", srvwrapper.New(hand.Handle))

	err = http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
