package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	createorder := &createorder.Handler{}
	http.Handle("/createOrder", srvwrapper.New(createorder.Handle))

	listorder := &listorder.Handler{}
	http.Handle("/listOrder", srvwrapper.New(listorder.Handle))

	orderpayed := &orderpayed.Handler{}
	http.Handle("/orderPayed", srvwrapper.New(orderpayed.Handle))

	cancelorder := &cancelorder.Handler{}
	http.Handle("/cancelOrder", srvwrapper.New(cancelorder.Handle))

	stocks := &stocks.Handler{}
	http.Handle("/stocks", srvwrapper.New(stocks.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR main ListenAndServe:", err)
	}
}
