package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/libs/kafka"
	"route256/libs/transactor"
	lomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	cancelOrder "route256/loms/internal/daemons/cancel-order"
	sendOrder "route256/loms/internal/daemons/send-order"
	"route256/loms/internal/domain"
	repository "route256/loms/internal/repository/postgres"
	"route256/loms/internal/sender"
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
		config.ConfigData.Services.LomsPgBouncer.UserDB,
		config.ConfigData.Services.LomsPgBouncer.PasswordDB,
		config.ConfigData.Services.LomsPgBouncer.Host,
		config.ConfigData.Services.LomsPgBouncer.Port,
		config.ConfigData.Services.LomsPgBouncer.NameDB)

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
	repo := repository.NewRepository(queryEngineProvider)

	domainRepository := domain.NewRepository(repo, queryEngineProvider)

	businessLogic := domain.NewService(domainRepository)

	desc.RegisterLOMSV1Server(s, lomsV1.NewLomsV1(businessLogic))

	// запускаем фоном отмену заказов
	cancelOrderDaemon := cancelOrder.NewCancelOrderDaemon(businessLogic)
	go cancelOrderDaemon.RunCancelDaemon(
		config.ConfigData.Services.CancelOrderDaemon.WorkersCount,
		time.Minute*time.Duration(config.ConfigData.Services.CancelOrderDaemon.CancelOrderTimeInMinutes))

	// запускаем фоном отмену заказов
	var brokers = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}
	producer, err := kafka.NewSyncProducer(brokers)
	if err != nil {
		log.Fatalln(err)
	}

	orderSender := sender.NewOrderSender(
		producer,
		"orders",
	)

	sendOrderDaemon := sendOrder.NewSendOrderDaemon(businessLogic, orderSender)
	go sendOrderDaemon.RunSendDaemon(
		5,
		"orders")

	log.Println("grpc server at", config.ConfigData.Services.Loms.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", err)
	}
}
