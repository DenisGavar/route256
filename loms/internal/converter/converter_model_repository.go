package converter

import (
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"
)

func FromRepositoryToMolelStocksResponse(stockItems []schema.StockItem) *model.StocksResponse {
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

func FromRepositoryToMolelStockItem(stockItem schema.StockItem) *model.StockItem {
	return &model.StockItem{
		WarehouseId: stockItem.WarehouseId,
		Count:       stockItem.Count,
	}
}
