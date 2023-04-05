package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "requests_total",
	},
		[]string{"full_method"},
	)
	HistogramResponseTimeServer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "histogram_response_time_server_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"full_method", "status"},
	)
	HistogramResponseTimeClient = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "histogram_response_time_client_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"full_method", "status"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	RequestsCounter.WithLabelValues(info.FullMethod).Inc()

	timeStart := time.Now()

	resp, err := handler(ctx, req)

	elapsed := time.Since(timeStart)

	statusGRPC := status.Code(err)

	HistogramResponseTimeServer.WithLabelValues(info.FullMethod, statusGRPC.String()).Observe(elapsed.Seconds())

	return resp, err
}

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	timeStart := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	elapsed := time.Since(timeStart)

	statusGRPC := status.Code(err)

	HistogramResponseTimeClient.WithLabelValues(method, statusGRPC.String()).Observe(elapsed.Seconds())

	return err
}
