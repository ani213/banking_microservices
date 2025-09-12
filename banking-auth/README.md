# For migration up Run
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" up

# For migration down Run
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" down