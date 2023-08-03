package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	hand := &stocks.Handler{}
	http.Handle("/stocks", srvwrapper.New(hand.Handle))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
