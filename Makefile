include .env
export

compose:
	@echo "Starting docker-compose..."
	@docker-compose up --build

migrate/up:
	@echo "Migrating up..."
	@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ./db/migrations up

migrate/down:
	@echo "Migrating down..."
	@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ./db/migrations down

sqlc:
	@echo "Generating sqlc..."
	@sqlc generate

test:
	@echo "Running tests..."
	@go test -count=1 -v ./cmd

.PHONY: compose migrate/up migrate/down sqlc test
.SILENT: compose migrate/up migrate/down sqlc test