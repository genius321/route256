package main

import (
	"encoding/json"
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

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (h addToCartHandle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	requestIn := Request{}

	token := request.Header.Get("Authorization")

	log.Printf("Token: %s", token)

	err := json.NewDecoder(request.Body).Decode(&requestIn)
	if err != nil {
		log.Printf("ERR: parse body: %s", err)
	}
	// data, err := io.ReadAll(request.Body)
	// if err != nil {
	// 	log.Printf("ERR: parse body: %s", err)
	// }
	// err = json.Unmarshal(data, &requestIn)
	// if err != nil {
	// 	log.Printf("ERR: parse body: %s", err)
	// }
	log.Printf("%+v", requestIn)

	if requestIn.User == 0 {
		response.WriteHeader(http.StatusBadRequest)
		_, _ = response.Write([]byte("User not set!\n"))
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write([]byte("{\"status\": \"done\"}\n"))
}
