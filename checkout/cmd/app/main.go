package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/grpc/loms"
	productService "route256/checkout/internal/clients/grpc/product-service"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/model"
	repository "route256/checkout/internal/repository/postgres"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/limiter"
	"route256/libs/logger"
	"route256/libs/transactor"
	workerPool "route256/libs/worker-pool"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
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

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := runGRPC()
		if err != nil {
			logger.Fatal("running GRPC", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := runHTTP(ctx)
		if err != nil {
			logger.Fatal("running HTTP", zap.Error(err))
		}
	}()

	wg.Wait()
}

func runGRPC() error {
	lis, err := net.Listen("tcp", config.ConfigData.Services.Checkout.GRPCPort)
	if err != nil {
		logger.Error("failed to listen grpc", zap.String("port", config.ConfigData.Services.Checkout.GRPCPort))
		return err
	}

	s := grpc.NewServer()
	reflection.Register(s)

	// создаём клиентов
	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to creating client loms", zap.String("address", config.ConfigData.Services.Loms.Address))
		return err
	}
	defer lomsConn.Close()
	lomsClient := loms.New(lomsConn)

	productServiceConn, err := grpc.Dial(config.ConfigData.Services.ProductService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to creating client product service", zap.String("address", config.ConfigData.Services.ProductService.Address))
		return err
	}
	defer productServiceConn.Close()
	productServiceClient := productService.New(productServiceConn, config.ConfigData.Services.ProductService.Token)

	rateLimit := config.ConfigData.Services.ProductService.RateLimit

	productServiceSettings := domain.NewProductServiceSettings(
		limiter.NewLimiter(time.Second, rateLimit),
	)

	productService := domain.NewProductService(productServiceClient, productServiceSettings)

	// подключаемся к БД
	ctx, cacnel := context.WithCancel(context.Background())
	defer cacnel()

	psqlConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.ConfigData.Services.CheckoutPgBouncer.UserDB,
		config.ConfigData.Services.CheckoutPgBouncer.PasswordDB,
		config.ConfigData.Services.CheckoutPgBouncer.Host,
		config.ConfigData.Services.CheckoutPgBouncer.Port,
		config.ConfigData.Services.CheckoutPgBouncer.NameDB)

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

	// создаём worker pool для обработки запросов list cart из product service
	wp := workerPool.New[*model.CartItem, *model.CartItem](
		ctx,
		config.ConfigData.Services.ProductService.ListCartWorkersCount,
	)
	wp.Init(ctx)

	businessLogic := domain.NewService(lomsClient, productService, domainRepository, wp)

	desc.RegisterCheckoutV1Server(s, checkoutV1.NewCheckoutV1(businessLogic))

	logger.Info("grpc server at", zap.String("port", config.ConfigData.Services.Checkout.GRPCPort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve grpc server")
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

	logger.Info("http served at", zap.String("port", config.ConfigData.Services.Checkout.HTTPPort))

	return http.ListenAndServe(config.ConfigData.Services.Checkout.HTTPPort, mux)
}
