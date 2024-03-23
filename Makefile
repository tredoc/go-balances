include .env

compose:
	@echo "Starting docker-compose..."
	@docker-compose up --build

migrate/up:
	@echo "Migrating up..."
	@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ./db/migrations up

migrate/down:
	@echo "Migrating down..."
	@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ./db/migrations down

test:
	@echo "Running tests..."
	@go test -v -cover ./...

.PHONY: compose migrate/up migrate/down test
.SILENT: compose migrate/up migrate/down test