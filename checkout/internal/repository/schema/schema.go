package schema

type CartItem struct {
	BasketId int64  `db:"id"`
	Sku      uint32 `db:"sku"`
	Count    uint32 `db:"count"`
}
