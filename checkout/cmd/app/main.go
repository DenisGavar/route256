package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	addtocart "route256/checkout/internal/handlers/add-to-cart"
	deletefromcart "route256/checkout/internal/handlers/delete-from-cart"
	listcart "route256/checkout/internal/handlers/list-cart"
	"route256/checkout/internal/handlers/purchase"
	srvwrapper "route256/libs/server-wrapper"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := loms.New(config.ConfigData.Services.Loms.Url)
	productServiseClient := productservice.New(config.ConfigData.Services.ProductService.Url, config.ConfigData.Services.ProductService.Token)

	busineddLogic := domain.New(lomsClient, lomsClient, productServiseClient)

	addToCartHandler := addtocart.New(busineddLogic)
	deleteFromCartHandler := deletefromcart.New()
	listCartHandler := listcart.New(busineddLogic)
	purchaseHandler := purchase.New(busineddLogic)

	http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))
	http.Handle("/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))
	http.Handle("/listCart", srvwrapper.New(listCartHandler.Handle))
	http.Handle("/purchase", srvwrapper.New(purchaseHandler.Handle))

	log.Println("listening http at", config.ConfigData.Services.Checkout.Port)
	err = http.ListenAndServe(config.ConfigData.Services.Checkout.Port, nil)
	log.Fatal("cannot listen http", err)
}
