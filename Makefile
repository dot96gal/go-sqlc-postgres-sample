include .env

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	DOCKER_HOST=${TEST_DOCKER_HOST} \
	TEST_POSTGRES_DB=${TEST_POSTGRES_DB} \
	TEST_POSTGRES_USER=${TEST_POSTGRES_USER} \
	TEST_POSTGRES_PASSWORD=${TEST_POSTGRES_PASSWORD} \
	go test -race ./...

.PHONY: dev
dev:
	POSTGRES_DB=${POSTGRES_DB} \
	POSTGRES_USER=${POSTGRES_USER} \
	POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	go run ./...

.PHONY: docker-compose-up
docker-compose-up:
	docker compose -f ./docker/compose.yml up -d

.PHONY: docker-compose-down
docker-compose-down:
	docker compose -f ./docker/compose.yml down

.PHONY: psql
psql:
	psql -h localhost -p 5432 -U ${POSTGRES_USER} -d ${POSTGRES_DB}

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: golang-migrate-up
golang-migrate-up:
	migrate -path db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" --verbose up
	
.PHONY: golang-migrate-down
golang-migrate-down:
	migrate -path db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" --verbose down
	
.PHONY: golang-migrate-drop
golang-migrate-drop:
	migrate -path db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" --verbose drop
