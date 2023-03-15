package model

type AddToCartRequest struct {
	// user ID
	User int64
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
}

type DeleteFromCartRequest struct {
	// user ID
	User int64
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
}

type ListCartRequest struct {
	// user ID
	User int64
}

type CartItem struct {
	// stock keeping unit - единица складского учёта
	Sku   uint32
	Count uint32
	// наименование товара
	Name string
	// цена товара
	Price uint32
	// ID записи в БД
	CartId int64
}

type ListCartResponse struct {
	Items      []*CartItem
	TotalPrice uint32
}

type PurchaseRequest struct {
	// user ID
	User int64
}

type PurchaseResponse struct {
	OrderId int64
}
