package loms

//go:generate mockgen -source="client.go" -destination="mocks/client_mock.go" -package=mocks . LomsClient
