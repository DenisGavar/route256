package converter

import (
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"
)

func FromRepositoryToMolelStocksResponse(stockItems []*schema.StockItem) *model.StocksResponse {
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
	if order == nil {
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

func FromRepositoryToMolelReserves(reserveItems []*schema.ReserveItem) *model.Reserve {
	items := make([]*model.ReserveStocksItem, 0, len(reserveItems))
	for _, i := range reserveItems {
		items = append(items, FromRepositoryToMolelReserveItem(i))
	}

	return &model.Reserve{
		ReserveItems: items,
	}
}

func FromRepositoryToMolelReserveItem(reserveItem *schema.ReserveItem) *model.ReserveStocksItem {
	return &model.ReserveStocksItem{
		WarehouseId: reserveItem.WarehouseId,
		Sku:         reserveItem.Sku,
		Count:       reserveItem.Count,
	}
}

func FromRepositoryToMolelCancelOrderRequestSlice(cancelOrderRequest []*schema.CancelOrderRequest) []*model.CancelOrderRequest {
	cancelOrderRequestSlice := make([]*model.CancelOrderRequest, 0, len(cancelOrderRequest))
	for _, i := range cancelOrderRequest {
		cancelOrderRequestSlice = append(cancelOrderRequestSlice, FromRepositoryToMolelCancelOrderRequest(i))
	}

	return cancelOrderRequestSlice
}

func FromRepositoryToMolelCancelOrderRequest(cancelOrderRequest *schema.CancelOrderRequest) *model.CancelOrderRequest {
	return &model.CancelOrderRequest{
		OrderId: cancelOrderRequest.OrderId,
	}
}

func FromRepositoryToMolelOrderMessageSlice(orderMessages []*schema.OrderMessage) []*model.OrderMessage {
	orderMessagesSlice := make([]*model.OrderMessage, 0, len(orderMessages))
	for _, i := range orderMessages {
		orderMessagesSlice = append(orderMessagesSlice, FromRepositoryToMolelOrderMessage(i))
	}

	return orderMessagesSlice
}

func FromRepositoryToMolelOrderMessage(orderMessage *schema.OrderMessage) *model.OrderMessage {
	return &model.OrderMessage{
		Id:        orderMessage.Id,
		OrderId:   orderMessage.OrderId,
		Message:   orderMessage.Message,
		CreatedAt: orderMessage.CreatedAt,
	}
}
