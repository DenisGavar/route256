package main

import (
	"log"
	"net"
	lomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lis, err := net.Listen("tcp", config.ConfigData.Services.Loms.Port)
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterLOMSV1Server(s, lomsV1.NewLomsV1())

	log.Println("grpc server at", config.ConfigData.Services.Loms.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}
