package product_service_cached

import (
	productServiceGRPCClient "route256/checkout/internal/clients/grpc/product-service"
	"route256/libs/cache"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "cache_hits_total",
	},
	)

	cacheErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "cache_errors_total",
	},
	)

	cacheRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "cache_requests_total",
	},
	)
)

type cachedClient struct {
	cache cache.Cache

	directClient productServiceGRPCClient.ProductServiceClient
}

var _ productServiceGRPCClient.ProductServiceClient = &cachedClient{}

func NewCachedClient(cache cache.Cache, directClient productServiceGRPCClient.ProductServiceClient) productServiceGRPCClient.ProductServiceClient {
	return &cachedClient{
		cache:        cache,
		directClient: directClient,
	}
}
