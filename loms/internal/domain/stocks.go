package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) Stocks(ctx context.Context, req *model.StocksRequest) (*model.StocksResponse, error) {
	// получаем остатки на складах
	logger.Debug("loms domain", zap.String("handler", "Stocks"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := s.repository.lomsRepository.Stocks(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingStocks.Error())
	}

	return response, nil
}
