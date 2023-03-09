package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/libs/transactor"
	lomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	repository "route256/loms/internal/repository/postgres"
	desc "route256/loms/pkg/loms_v1"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

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

	// подключаемся к БД
	ctx, cacnel := context.WithCancel(context.Background())
	defer cacnel()

	psqlConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.ConfigData.Services.LomsDB.User,
		config.ConfigData.Services.LomsDB.Password,
		config.ConfigData.Services.LomsDB.Host,
		config.ConfigData.Services.LomsDB.Port,
		config.ConfigData.Services.LomsDB.DBName)

	// пул соединений
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// настройки
	configDB := pool.Config()
	configDB.MaxConnIdleTime = time.Minute
	configDB.MaxConnLifetime = time.Hour
	configDB.MinConns = 2
	configDB.MaxConns = 10

	queryEngineProvider := transactor.NewTransactionManager(pool)
	repo := repository.NewRepo(queryEngineProvider)

	domainRepository := domain.NewRepository(repo, queryEngineProvider)

	businessLogic := domain.NewService(domainRepository)

	desc.RegisterLOMSV1Server(s, lomsV1.NewLomsV1(businessLogic))

	log.Println("grpc server at", config.ConfigData.Services.Loms.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}
