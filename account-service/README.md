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

# for docker account service step-1 create build
docker build -t account-service .
# for run account service step-2 run
docker run -d --name account-service -p 8081:8081   -e DATABASE_URL="postgres://postgre
s:password@host.docker.internal:5432/bank?sslmode=disable" -e SECRET_KEY=super-secret-
key -e EMAIL_SERVER=http://localhost:8083 account-service
