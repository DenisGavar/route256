package model

import "time"

type CancelOrderRequest struct {
	OrderId int64
}

type CreateOrderRequest struct {
	User    int64
	OrderId int64
	Items   []*OrderItem
}

type OrderItem struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type CreateOrderResponse struct {
	OrderId int64
}

type ListOrderRequest struct {
	OrderId int64
}

// статус заказа (new | awaiting payment | failed | payed | cancelled)

const (
	OrderStatusNew             string = "new"
	OrderStatusAwaitingPayment string = "awaiting payment"
	OrderStatusFailed          string = "failed"
	OrderStatusPayed           string = "payed"
	OrderStatusCancelled       string = "cancelled"
)

type ListOrderResponse struct {
	Status string       `json:"status"`
	User   int64        `json:"user"`
	Items  []*OrderItem `json:"items"`
}

type OrderPayedRequest struct {
	OrderId int64
}

type StocksRequest struct {
	Sku uint32
}

type StockItem struct {
	WarehouseId int64
	Count       uint64
}

type StocksResponse struct {
	Stocks []*StockItem
}

type ReserveStocksItem struct {
	WarehouseId int64
	Sku         uint32
	Count       uint64
}

type Reserve struct {
	ReserveItems []*ReserveStocksItem
}

type OrderMessage struct {
	Id        int64
	OrderId   int64
	Message   string
	CreatedAt time.Time
}
