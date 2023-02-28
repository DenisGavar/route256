package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	productService "route256/checkout/internal/clients/product-service"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	addToCart "route256/checkout/internal/handlers/add-to-cart"
	deleteFromCart "route256/checkout/internal/handlers/delete-from-cart"
	listCart "route256/checkout/internal/handlers/list-cart"
	"route256/checkout/internal/handlers/purchase"
	serverWrapper "route256/libs/server-wrapper"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := loms.New(config.ConfigData.Services.Loms.Url)
	productServiseClient := productService.New(config.ConfigData.Services.ProductService.Url, config.ConfigData.Services.ProductService.Token)

	businessLogic := domain.New(lomsClient, lomsClient, productServiseClient)

	addToCartHandler := addToCart.New(businessLogic)
	deleteFromCartHandler := deleteFromCart.New()
	listCartHandler := listCart.New(businessLogic)
	purchaseHandler := purchase.New(businessLogic)

	http.Handle("/addToCart", serverWrapper.New(addToCartHandler.Handle))
	http.Handle("/deleteFromCart", serverWrapper.New(deleteFromCartHandler.Handle))
	http.Handle("/listCart", serverWrapper.New(listCartHandler.Handle))
	http.Handle("/purchase", serverWrapper.New(purchaseHandler.Handle))

	log.Println("listening http at", config.ConfigData.Services.Checkout.Port)
	err = http.ListenAndServe(config.ConfigData.Services.Checkout.Port, nil)
	log.Fatal("cannot listen http", err)
}
