package repository

import (
	"route256/libs/transactor"
)

type repo struct {
	queryEngineProvider transactor.QueryEngineProvider
}

func NewRepo(queryEngineProvider transactor.QueryEngineProvider) *repo {
	return &repo{
		queryEngineProvider: queryEngineProvider,
	}
}
