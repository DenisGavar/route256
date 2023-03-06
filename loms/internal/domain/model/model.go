package model

type CancelOrderRequest struct {
	OrderId int64
}

type CreateOrderRequest struct {
	User  int64
	Items []*OrderItem
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
type OderStatus int32

const (
	OderStatus_new              OderStatus = 0
	OderStatus_awaiting_payment OderStatus = 1
	OderStatus_failed           OderStatus = 2
	OderStatus_payed            OderStatus = 3
	OderStatus_cancelled        OderStatus = 4
)

type ListOrderResponse struct {
	Status OderStatus
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
