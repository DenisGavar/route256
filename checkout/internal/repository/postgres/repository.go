package repository

import (
	"route256/libs/transactor"
)

const (
	basketsTable = "baskets"
)

type repo struct {
	queryEngineProvider transactor.QueryEngineProvider
}

func NewRepo(queryEngineProvider transactor.QueryEngineProvider) *repo {
	return &repo{
		queryEngineProvider: queryEngineProvider,
	}
}
