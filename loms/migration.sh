goose -dir ./migrations postgres "postgres://user:password@localhost:5432/loms?sslmode=disable" status

goose -dir ./migrations postgres "postgres://user:password@localhost:5432/loms?sslmode=disable" up