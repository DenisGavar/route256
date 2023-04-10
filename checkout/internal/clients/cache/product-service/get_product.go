package product_service_cached

import (
	"context"
	"fmt"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/logger"
	"strconv"

	"go.uber.org/zap"
)

func (c *cachedClient) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	logger.Debug("product-service cached client", zap.String("handler", "GetProduct"), zap.String("request", fmt.Sprintf("%+v", req)))

	res, ok := c.getFromCacheMetered(ctx, req)
	if !ok {
		return c.requestAndCache(ctx, req)
	}

	// если нашли в кэше
	return res, nil
}

func (c *cachedClient) requestAndCache(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	res, err := c.directClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	go func() {
		if ok := c.cache.Set(context.Background(), strconv.Itoa(int(req.Sku)), res); !ok {
			cacheErrorsTotal.Inc()

			return
		}
	}()

	return res, nil
}

func (c *cachedClient) getFromCacheMetered(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, bool) {
	value, ok := c.getFromCache(ctx, req)
	cacheRequestsTotal.Inc()

	if ok {
		cacheHitsTotal.Inc()
	}

	return value, ok
}

func (c *cachedClient) getFromCache(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, bool) {
	res, ok := c.cache.Get(ctx, strconv.Itoa(int(req.Sku))).(*product.GetProductResponse)
	if !ok {
		cacheErrorsTotal.Inc()
		return nil, false
	}

	return res, true
}
