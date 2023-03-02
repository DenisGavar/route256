package main

import (
	"log"
	"net"
	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/grpc/loms"
	productService "route256/checkout/internal/clients/grpc/product-service"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// создаём клиентов
	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to creating client loms", err)
	}
	defer lomsConn.Close()
	lomsClient := loms.New(lomsConn)

	productServiceConn, err := grpc.Dial(config.ConfigData.Services.ProductService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to creating client loms", err)
	}
	defer productServiceConn.Close()
	productServiceClient := productService.New(productServiceConn, config.ConfigData.Services.ProductService.Token)

	businessLogic := domain.NewModel(lomsClient, productServiceClient)

	desc.RegisterCheckoutV1Server(s, checkoutV1.NewCheckoutV1(businessLogic))

	log.Println("grpc server at", config.ConfigData.Services.Checkout.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}
