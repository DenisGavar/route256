package main

import (
	"context"
	"log"
	"net"
	"net/http"
	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/grpc/loms"
	productService "route256/checkout/internal/clients/grpc/product-service"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := runGRPC()
		if err != nil {
			log.Fatal("running GRPC", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := runHTTP(ctx)
		if err != nil {
			log.Fatal("running HTTP", err)
		}
	}()

	wg.Wait()
}

func runGRPC() error {
	lis, err := net.Listen("tcp", config.ConfigData.Services.Checkout.GRPCPort)
	if err != nil {
		log.Printf("failed to listen")
		return err
	}

	s := grpc.NewServer()
	reflection.Register(s)

	// создаём клиентов
	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to creating client loms")
		return err
	}
	defer lomsConn.Close()
	lomsClient := loms.New(lomsConn)

	productServiceConn, err := grpc.Dial(config.ConfigData.Services.ProductService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to creating client product service")
		return err
	}
	defer productServiceConn.Close()
	productServiceClient := productService.New(productServiceConn, config.ConfigData.Services.ProductService.Token)

	businessLogic := domain.NewModel(lomsClient, productServiceClient)

	desc.RegisterCheckoutV1Server(s, checkoutV1.NewCheckoutV1(businessLogic))

	log.Println("grpc server at", config.ConfigData.Services.Checkout.GRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve")
		return err
	}

	return nil
}

func runHTTP(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := desc.RegisterCheckoutV1HandlerFromEndpoint(ctx, mux, config.ConfigData.Services.Checkout.GRPCPort, opts)
	if err != nil {
		return err
	}

	log.Println("http served at", config.ConfigData.Services.Checkout.HTTPPort)

	return http.ListenAndServe(config.ConfigData.Services.Checkout.HTTPPort, mux)

}
