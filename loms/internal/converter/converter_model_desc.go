package converter

import (
	"route256/loms/internal/domain/model"
	desc "route256/loms/pkg/loms_v1"
)

func FromDescToMolelCreateOrderRequest(createOrderRequest *desc.CreateOrderRequest) *model.CreateOrderRequest {
	if createOrderRequest == nil {
		return nil
	}

	items := make([]*model.OrderItem, 0, len(createOrderRequest.GetItems()))
	for _, i := range createOrderRequest.GetItems() {
		items = append(items, FromDescToMolelOrderItem(i))
	}

	return &model.CreateOrderRequest{
		User:  createOrderRequest.GetUser(),
		Items: items,
	}
}

func FromDescToMolelOrderItem(orderItem *desc.OrderItem) *model.OrderItem {

	return &model.OrderItem{
		Sku:   orderItem.GetSku(),
		Count: uint16(orderItem.GetCount()),
	}
}

func FromModelToDescCreateOrderResponse(createOrderResponse *model.CreateOrderResponse) *desc.CreateOrderResponse {
	if createOrderResponse == nil {
		return nil
	}

	return &desc.CreateOrderResponse{
		OrderId: createOrderResponse.OrderId,
	}
}

func FromDescToMolelListOrderRequest(listOrderRequest *desc.ListOrderRequest) *model.ListOrderRequest {
	if listOrderRequest == nil {
		return nil
	}

	return &model.ListOrderRequest{
		OrderId: listOrderRequest.GetOrderId(),
	}
}

func FromModelToDescListOrderResponse(listOrderResponse *model.ListOrderResponse) *desc.ListOrderResponse {
	if listOrderResponse == nil {
		return nil
	}

	items := make([]*desc.OrderItem, 0, len(listOrderResponse.Items))
	for _, i := range listOrderResponse.Items {
		items = append(items, FromMolelToDescOrderItem(i))
	}

	return &desc.ListOrderResponse{
		Status: listOrderResponse.Status,
		User:   listOrderResponse.User,
		Items:  items,
	}
}

func FromMolelToDescOrderItem(orderItem *model.OrderItem) *desc.OrderItem {
	return &desc.OrderItem{
		Sku:   orderItem.Sku,
		Count: uint32(orderItem.Count),
	}
}

func FromDescToMolelOrderPayedRequest(orderPayedRequest *desc.OrderPayedRequest) *model.OrderPayedRequest {
	if orderPayedRequest == nil {
		return nil
	}

	return &model.OrderPayedRequest{
		OrderId: orderPayedRequest.GetOrderId(),
	}
}

func FromDescToMolelCancelOrderRequest(cancelOrderRequest *desc.CancelOrderRequest) *model.CancelOrderRequest {
	if cancelOrderRequest == nil {
		return nil
	}

	return &model.CancelOrderRequest{
		OrderId: cancelOrderRequest.GetOrderId(),
	}
}

func FromDescToMolelStocksRequest(stocksRequest *desc.StocksRequest) *model.StocksRequest {
	if stocksRequest == nil {
		return nil
	}

	return &model.StocksRequest{
		Sku: stocksRequest.GetSku(),
	}
}

func FromModelToDescStocksResponse(stocksResponse *model.StocksResponse) *desc.StocksResponse {
	if stocksResponse == nil {
		return nil
	}

	items := make([]*desc.StockItem, 0, len(stocksResponse.Stocks))
	for _, i := range stocksResponse.Stocks {
		items = append(items, FromMolelToDescStockItem(i))
	}

	return &desc.StocksResponse{
		Stocks: items,
	}
}

func FromMolelToDescStockItem(stockItem *model.StockItem) *desc.StockItem {
	return &desc.StockItem{
		WarehouseId: stockItem.WarehouseId,
		Count:       stockItem.Count,
	}
}
