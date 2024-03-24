include .env

# --- local ---

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
	POSTGRES_HOST=${POSTGRES_HOST} \
	POSTGRES_PORT=${POSTGRES_PORT} \
	POSTGRES_SSL_MODE=${POSTGRES_SSL_MODE} \
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


# --- supabase ---

.PHONY: dev-supabase
dev-supabase:
	POSTGRES_DB=${SUPABASE_DB} \
	POSTGRES_USER=${SUPABASE_USER} \
	POSTGRES_PASSWORD=${SUPABASE_PASSWORD} \
	POSTGRES_HOST=${SUPABASE_HOST} \
	POSTGRES_PORT=${SUPABASE_PORT} \
	POSTGRES_SSL_MODE=${SUPABASE_SSL_MODE} \
	go run ./...

.PHONY: psql-supabase
psql-supabase:
	psql -h ${SUPABASE_HOST} -p 5432 -d postgres -U ${SUPABASE_USER}

.PHONY: golang-migrate-up-supabase
golang-migrate-up-supabase:
	migrate -path db/migrations -database "postgres://${SUPABASE_USER}:${SUPABASE_PASSWORD}@${SUPABASE_HOST}:${SUPABASE_PORT}/${SUPABASE_DB}?sslmode=${SUPABASE_SSL_MODE}" --verbose up
	
.PHONY: golang-migrate-down-supabase
golang-migrate-down-supabase:
	migrate -path db/migrations -database "postgres://${SUPABASE_USER}:${SUPABASE_PASSWORD}@${SUPABASE_HOST}:${SUPABASE_PORT}/${SUPABASE_DB}?sslmode=${SUPABASE_SSL_MODE}" --verbose down
	
.PHONY: golang-migrate-drop-supabase
golang-migrate-drop-supabase:
	migrate -path db/migrations -database "postgres://${SUPABASE_USER}:${SUPABASE_PASSWORD}@${SUPABASE_HOST}:${SUPABASE_PORT}/${SUPABASE_DB}?sslmode=${SUPABASE_SSL_MODE}" --verbose drop
