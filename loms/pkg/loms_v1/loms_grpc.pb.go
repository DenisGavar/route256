// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: loms.proto

package loms_v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LOMSV1Client is the client API for LOMSV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LOMSV1Client interface {
	// Создает новый заказ для пользователя из списка переданных товаров. Товары при этом нужно зарезервировать на складе
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	// Показывает информацию по заказу
	ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error)
	// Помечает заказ оплаченным. Зарезервированные товары должны перейти в статус купленных
	OrderPayed(ctx context.Context, in *OrderPayedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Отменяет заказ, снимает резерв со всех товаров в заказе
	CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Возвращает количество товаров, которые можно купить с разных складов. Если товар был зарезерванован у кого-то в заказе и ждет оплаты, его купить нельзя
	Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error)
}

type lOMSV1Client struct {
	cc grpc.ClientConnInterface
}

func NewLOMSV1Client(cc grpc.ClientConnInterface) LOMSV1Client {
	return &lOMSV1Client{cc}
}

func (c *lOMSV1Client) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/loms_v1.LOMSV1/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSV1Client) ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error) {
	out := new(ListOrderResponse)
	err := c.cc.Invoke(ctx, "/loms_v1.LOMSV1/ListOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSV1Client) OrderPayed(ctx context.Context, in *OrderPayedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms_v1.LOMSV1/OrderPayed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSV1Client) CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms_v1.LOMSV1/CancelOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSV1Client) Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error) {
	out := new(StocksResponse)
	err := c.cc.Invoke(ctx, "/loms_v1.LOMSV1/Stocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LOMSV1Server is the server API for LOMSV1 service.
// All implementations must embed UnimplementedLOMSV1Server
// for forward compatibility
type LOMSV1Server interface {
	// Создает новый заказ для пользователя из списка переданных товаров. Товары при этом нужно зарезервировать на складе
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	// Показывает информацию по заказу
	ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error)
	// Помечает заказ оплаченным. Зарезервированные товары должны перейти в статус купленных
	OrderPayed(context.Context, *OrderPayedRequest) (*emptypb.Empty, error)
	// Отменяет заказ, снимает резерв со всех товаров в заказе
	CancelOrder(context.Context, *CancelOrderRequest) (*emptypb.Empty, error)
	// Возвращает количество товаров, которые можно купить с разных складов. Если товар был зарезерванован у кого-то в заказе и ждет оплаты, его купить нельзя
	Stocks(context.Context, *StocksRequest) (*StocksResponse, error)
	mustEmbedUnimplementedLOMSV1Server()
}

// UnimplementedLOMSV1Server must be embedded to have forward compatible implementations.
type UnimplementedLOMSV1Server struct {
}

func (UnimplementedLOMSV1Server) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedLOMSV1Server) ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrder not implemented")
}
func (UnimplementedLOMSV1Server) OrderPayed(context.Context, *OrderPayedRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderPayed not implemented")
}
func (UnimplementedLOMSV1Server) CancelOrder(context.Context, *CancelOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedLOMSV1Server) Stocks(context.Context, *StocksRequest) (*StocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stocks not implemented")
}
func (UnimplementedLOMSV1Server) mustEmbedUnimplementedLOMSV1Server() {}

// UnsafeLOMSV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LOMSV1Server will
// result in compilation errors.
type UnsafeLOMSV1Server interface {
	mustEmbedUnimplementedLOMSV1Server()
}

func RegisterLOMSV1Server(s grpc.ServiceRegistrar, srv LOMSV1Server) {
	s.RegisterService(&LOMSV1_ServiceDesc, srv)
}

func _LOMSV1_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSV1Server).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LOMSV1/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSV1Server).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMSV1_ListOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSV1Server).ListOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LOMSV1/ListOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSV1Server).ListOrder(ctx, req.(*ListOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMSV1_OrderPayed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPayedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSV1Server).OrderPayed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LOMSV1/OrderPayed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSV1Server).OrderPayed(ctx, req.(*OrderPayedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMSV1_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSV1Server).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LOMSV1/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSV1Server).CancelOrder(ctx, req.(*CancelOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMSV1_Stocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSV1Server).Stocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LOMSV1/Stocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSV1Server).Stocks(ctx, req.(*StocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LOMSV1_ServiceDesc is the grpc.ServiceDesc for LOMSV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LOMSV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loms_v1.LOMSV1",
	HandlerType: (*LOMSV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _LOMSV1_CreateOrder_Handler,
		},
		{
			MethodName: "ListOrder",
			Handler:    _LOMSV1_ListOrder_Handler,
		},
		{
			MethodName: "OrderPayed",
			Handler:    _LOMSV1_OrderPayed_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _LOMSV1_CancelOrder_Handler,
		},
		{
			MethodName: "Stocks",
			Handler:    _LOMSV1_Stocks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loms.proto",
}
