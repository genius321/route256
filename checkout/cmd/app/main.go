package main

import (
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	hand := addToCartHandle{}
	http.Handle("/addToCart", hand)

	err := http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
