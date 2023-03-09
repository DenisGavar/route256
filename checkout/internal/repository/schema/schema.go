package schema

type CartItem struct {
	Sku   uint32 `db:"sku"`
	Count uint32 `db:"count"`
}
