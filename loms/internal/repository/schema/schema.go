package schema

import "time"

type StockItem struct {
	StockId     int64  `db:"id"`
	Sku         int64  `db:"sku"`
	WarehouseId int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}

type Order struct {
	OrderId int64  `db:"id"`
	User    int64  `db:"user_id"`
	Status  string `db:"status"`
}

type OrderItem struct {
	Sku   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}

type ReserveItem struct {
	WarehouseId int64  `db:"warehouse_id"`
	Sku         uint32 `db:"sku"`
	Count       uint64 `db:"count"`
}

type CancelOrderRequest struct {
	OrderId int64 `db:"id"`
}

type OrderMessage struct {
	Id        int64     `db:"id"`
	OrderId   int64     `db:"orders_id"`
	Message   string    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
}
