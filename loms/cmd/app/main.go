package main

import (
	"context"
	"fmt"
	"net"
	"route256/libs/kafka"
	"route256/libs/logger"
	"route256/libs/tracing"
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

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger.Init()

	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	tracing.Init(logger.GetLogger(), config.ConfigData.Services.Loms.Name)

	lis, err := net.Listen("tcp", config.ConfigData.Services.Loms.Port)
	if err != nil {
		logger.Fatal("failed to listen grpc", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
		),
	)
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
		logger.Fatal("failed to creating pgxpool connection", zap.Error(err))
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

	// создаём producer для kafka
	producer, err := kafka.NewSyncProducer(config.ConfigData.Services.Kafka.Brokers)
	if err != nil {
		logger.Fatal("failed to creating kafka producer", zap.Error(err))
	}

	orderSender := sender.NewOrderSender(
		producer,
		config.ConfigData.Services.Kafka.TopicForOrders,
	)

	// запускаем фоном отправку смены статусов заказов в kafka
	sendOrderDaemon := sendOrder.NewSendOrderDaemon(businessLogic, orderSender)
	go sendOrderDaemon.RunSendDaemon(
		config.ConfigData.Services.Kafka.WorkersCount,
		config.ConfigData.Services.Kafka.TopicForOrders)

	logger.Info("grpc server at", zap.String("port", config.ConfigData.Services.Loms.Port))

	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve grpc server", zap.Error(err))
	}
}
