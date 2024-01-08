package grpc_wrapper

import (
	"route256/libs/metrics"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConn interface {
	Close() error
	GetClientConn() *grpc.ClientConn
}

type clientConn struct {
	grpcClientConn *grpc.ClientConn
}

func Dial(target string) (ClientConn, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
			metrics.UnaryClientInterceptor,
		),
	)
	if err != nil {
		return nil, err
	}

	return &clientConn{
		grpcClientConn: conn,
	}, nil
}

func (cc *clientConn) Close() error {
	return cc.grpcClientConn.Close()
}

func (cc *clientConn) GetClientConn() *grpc.ClientConn {
	return cc.grpcClientConn
}
