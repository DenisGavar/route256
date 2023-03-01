package main

import (
	"log"
	"net"
	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lis, err := net.Listen("tcp", config.ConfigData.Services.Checkout.Port)
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	businessLogic := domain.NewModel()

	desc.RegisterCheckoutV1Server(s, checkoutV1.NewCheckoutV1(businessLogic))

	log.Println("grpc server at", config.ConfigData.Services.Checkout.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}

	// lomsClient := loms.New(config.ConfigData.Services.Loms.Url)
	// productServiseClient := productservice.New(config.ConfigData.Services.ProductService.Url, config.ConfigData.Services.ProductService.Token)

	// busineddLogic := domain.New(lomsClient, lomsClient, productServiseClient)

	// addToCartHandler := addtocart.New(busineddLogic)
	// deleteFromCartHandler := deletefromcart.New()
	// listCartHandler := listcart.New(busineddLogic)
	// purchaseHandler := purchase.New(busineddLogic)

	// http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))
	// http.Handle("/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))
	// http.Handle("/listCart", srvwrapper.New(listCartHandler.Handle))
	// http.Handle("/purchase", srvwrapper.New(purchaseHandler.Handle))

	// log.Println("listening http at", config.ConfigData.Services.Checkout.Port)
	// err = http.ListenAndServe(config.ConfigData.Services.Checkout.Port, nil)
	// log.Fatal("cannot listen http", err)
}
