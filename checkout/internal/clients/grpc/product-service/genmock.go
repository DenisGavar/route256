package product_service

//go:generate mockgen -source="client.go" -destination="mocks/client_mock.go" -package=mocks . ProductServiceClient
