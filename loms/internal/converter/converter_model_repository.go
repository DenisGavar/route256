package converter

import (
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"
)

func FromRepositoryToMolelStocksResponse(stockItems []*schema.StockItem) *model.StocksResponse {
	if stockItems == nil {
		return nil
	}

	items := make([]*model.StockItem, 0, len(stockItems))
	for _, i := range stockItems {
		items = append(items, FromRepositoryToMolelStockItem(i))
	}

	return &model.StocksResponse{
		Stocks: items,
	}

}

func FromRepositoryToMolelStockItem(stockItem *schema.StockItem) *model.StockItem {
	return &model.StockItem{
		WarehouseId: stockItem.WarehouseId,
		Count:       stockItem.Count,
	}
}

func FromRepositoryToMolelListOrderResponse(order *schema.Order, orderItems []*schema.OrderItem) *model.ListOrderResponse {
	if order == nil || orderItems == nil {
		return nil
	}

	items := make([]*model.OrderItem, 0, len(orderItems))
	for _, i := range orderItems {
		items = append(items, FromRepositoryToMolelOrderItem(i))
	}

	return &model.ListOrderResponse{
		Status: order.Status,
		User:   order.User,
		Items:  items,
	}
}

func FromRepositoryToMolelOrderItem(orderItem *schema.OrderItem) *model.OrderItem {
	return &model.OrderItem{
		Sku:   orderItem.Sku,
		Count: orderItem.Count,
	}
}
