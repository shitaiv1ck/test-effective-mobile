include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d test-task-postgres

env-down:
	@docker compose down test-task-postgres

migrate-create:
	@if [ -z "${seq}"]; then \
		echo "pls, try again with seq=value"; \
		exit 1; \
	fi; \
	docker compose run --rm test-task-migrate \
	create -ext sql -dir migrations -seq ${seq}

migrate-up:
	@docker compose run --rm test-task-migrate \
	-path migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@test-task-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	up

migrate-down:
	@docker compose run --rm test-task-migrate \
	-path migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@test-task-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	down

migrate-force-1:
	@docker compose run --rm test-task-migrate \
	-path migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@test-task-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	force 1

app-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/test/main.go

app-deploy:
	@docker compose up -d --build test-task

app-deploy-stop:
	@docker compose down test-task

swagger-gen:
	@docker compose run --rm swagger \
	init \
	-g cmd/test/main.go \
	-o docs \
	--parseInternal \
	--parseDependency

