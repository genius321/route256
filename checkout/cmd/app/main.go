package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	err := config.Init()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	modelStockChecker := domain.NewModelStockChecker(loms.NewClientStocks(config.AppConfig.Services.Loms))
	modelProductGetter := domain.NewModelProductGetter(productservice.NewClientGetProduct(config.AppConfig.Services.ProductService))
	modelOrderCreater := domain.NewModelOrderCreater(loms.NewClientCreateOrder(config.AppConfig.Services.Loms))

	addtocart := &addtocart.Handler{
		Model: modelStockChecker,
	}
	http.Handle("/addToCart", srvwrapper.New(addtocart.Handle))

	deletefromcart := &deletefromcart.Handler{}
	http.Handle("/deleteFromCart", srvwrapper.New(deletefromcart.Handle))

	listcart := &listcart.Handler{
		Model: modelProductGetter,
	}
	http.Handle("/listCart", srvwrapper.New(listcart.Handle))

	purchase := &purchase.Handler{
		Model: modelOrderCreater,
	}
	http.Handle("/purchase", srvwrapper.New(purchase.Handle))

	err = http.ListenAndServe(port, nil)
	log.Fatalln("ERR: ", err)
}
