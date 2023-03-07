package schema

type StockItem struct {
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
