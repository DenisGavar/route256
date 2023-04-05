package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ozon",
		Subsystem: "http",
		Name:      "requests_total",
	},
		[]string{"full_method"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	RequestsCounter.WithLabelValues(info.FullMethod).Inc()

	resp, err := handler(ctx, req)

	return resp, err
}
