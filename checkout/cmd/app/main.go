package main

import (
	"log"
	"net/http"
)

const port = ":8080"

type addToCartHandle struct {
}

func main() {

	hand := addToCartHandle{}
	http.Handle("/addToCart", hand)

	err := http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}

func (h addToCartHandle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write([]byte("{\"status\": \"done\"}\n"))
}
