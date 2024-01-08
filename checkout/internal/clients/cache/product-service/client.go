package product_service_cached

import (
	"context"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/cache"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "hits_total",
	},
	)

	cacheErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "errors_total",
	},
	)

	cacheRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "requests_total",
	},
	)

	HistogramResponseTimeCache = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.000001, 2, 16),
	},
	)
)

type CachedClient interface {
	GetProduct(context.Context, *product.GetProductRequest) (*product.GetProductResponse, bool)
	SetToCache(ctx context.Context, req *product.GetProductRequest, res *product.GetProductResponse)
}

type cachedClient struct {
	cache cache.Cache
}

var _ CachedClient = &cachedClient{}

func NewCachedClient(cache cache.Cache) CachedClient {
	return &cachedClient{
		cache: cache,
	}
}
