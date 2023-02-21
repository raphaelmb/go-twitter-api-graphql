mock:
	mockery --all --keeptree

migrate:
	migrate -source file://postgres/migrations \
					-database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" up

rollback:
	migrate -source file://postgres/migrations \
					-database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" down

drop:
	migrate -source file://postgres/migrations \
					-database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" drop

migration:
	@read -p "Enter migration name: " name; \
			migrate create -ext sql -dir postgres/migrations $$name

run:
	go run cmd/graphqlserver/main.go

databaseup:
	docker compose up -d

databasedown:
	docker compose down

generate:
	go generate ./..