package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	stocks := &stocks.Handler{}
	http.Handle("/stocks", srvwrapper.New(stocks.Handle))

	createorder := &createorder.Handler{}
	http.Handle("/createOrder", srvwrapper.New(createorder.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
