package domain

import "context"

var _ Model = (*model)(nil)

type Model interface {
	AddToCart(ctx context.Context, req *AddToCartRequest) error
	ListCart(ctx context.Context, req *ListCartRequest) (*ListCartResponse, error)
	Purchase(ctx context.Context, req *PurchaseRequest) (*PurchaseResponse, error)
}

type StocksChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductGetter interface {
	GetProduct(ctx context.Context, sku uint32) (*Product, error)
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64) (*Order, error)
}

type model struct {
	//stocksChecker StocksChecker
	//productGetter ProductGetter
	//orderCreator  OrderCreator
}

//func NewModel(stocksChecker StocksChecker, orderCreator OrderCreator, productGetter ProductGetter) *model {
func NewModel() *model {
	return &model{
		//stocksChecker: stocksChecker,
		//orderCreator:  orderCreator,
		//productGetter: productGetter,
	}
}
