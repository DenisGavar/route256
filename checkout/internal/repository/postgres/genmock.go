package repository

//go:generate mockgen -source="repository.go" -destination="mocks/repository_mock.go" -package=mocks . CheckoutRepository
