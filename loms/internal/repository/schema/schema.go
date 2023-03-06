package schema

type StockItem struct {
	WarehouseId int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}
