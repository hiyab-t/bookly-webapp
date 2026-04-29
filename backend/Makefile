include .env


run/api:
	go run ./cmd/api ${ARGS} -cors-trusted-origins ${ORIGIN}
start/containers:
	@echo 'creating service containers'
	docker compose up -d
db/migrations/up: start/containers 
	sleep 3
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} up
.PHONY: run/api start/containers db/migrations/up

