# For migration up Run
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" up

# For migration down Run
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" down


# For Creating a migration file Run
migrate create -ext sql -dir migrations -seq create_role_table

migrations-> is migration directory
create_role_table-> is migration file name

# When getting error like version 3 or versin 2
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" force 3/2