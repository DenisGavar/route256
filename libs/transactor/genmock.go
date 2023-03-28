package transactor

//go:generate mockgen -source="transactor.go" -destination="mocks/db_mock.go" -package=mocks . DB
