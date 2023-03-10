package loms

import (
	lomsServiceAPI "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
)

type client struct {
	lomsClient lomsServiceAPI.LOMSV1Client
}

func New(cc *grpc.ClientConn) *client {
	return &client{
		lomsClient: lomsServiceAPI.NewLOMSV1Client(cc),
	}
}
