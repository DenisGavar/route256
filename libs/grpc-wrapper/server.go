package grpc_wrapper

import (
	"net"
	"route256/libs/metrics"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Serve(lis net.Listener) error
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
}

type server struct {
	grpcServer *grpc.Server
}

func NewServer() Server {
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
			metrics.UnaryServerInterceptor,
		),
	)
	reflection.Register(s)

	return &server{
		grpcServer: s,
	}
}

func (s *server) Serve(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

func (s *server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.grpcServer.RegisterService(desc, impl)
}
