package model

type CancelOrderRequest struct {
	OrderId int64
}

type CreateOrderRequest struct {
	User    int64
	OrderId int64
	Items   []*OrderItem
}

type OrderItem struct {
	Sku   uint32
	Count uint16
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
	Status string
	User   int64
	Items  []*OrderItem
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
	Part        bool
}

type Reserve struct {
	ReserveItems []*ReserveStocksItem
}
