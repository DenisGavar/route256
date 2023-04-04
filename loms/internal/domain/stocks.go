package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) Stocks(ctx context.Context, req *model.StocksRequest) (*model.StocksResponse, error) {
	// получаем остатки на складах
	logger.Debug("loms domain", zap.String("handler", "Stocks"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms domain Stocks processing")
	defer span.Finish()

	span.SetTag("sku", req.Sku)

	response, err := s.repository.lomsRepository.Stocks(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingStocks.Error())
	}

	return response, nil
}
