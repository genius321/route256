package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	err := config.Init()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	model := domain.New(loms.New(config.AppConfig.Services.Loms))

	addtocart := &addtocart.Handler{
		Model: model,
	}
	http.Handle("/addToCart", srvwrapper.New(addtocart.Handle))

	deletefromcart := &deletefromcart.Handler{}
	http.Handle("/deleteFromCart", srvwrapper.New(deletefromcart.Handle))

	err = http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
