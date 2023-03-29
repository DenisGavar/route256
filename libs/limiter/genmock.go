package limiter

//go:generate mockgen -source="limiter.go" -destination="mocks/limiter_mock.go" -package=mocks . Limiter
