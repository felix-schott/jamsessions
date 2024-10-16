include .env

.PHONY = test-ci test-db test-geocoding test-server test-cli dev-server migrations

TEST_DB_PORT = 54930
DEV_SERVER_PORT = 36230

DB_URL := "host=localhost port=$(TEST_DB_PORT) user=read_write password=$(READ_WRITE_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable"

internal/db/models.go: internal/db/schema.sql internal/db/query.sql internal/db/sqlc.yaml 
	sqlc generate -f internal/db/sqlc.yaml

test-ci: internal/db/models.go bin/dbcli # run in ci pipeline
	@TEST_DB_PORT=$(TEST_DB_PORT) docker compose up test_db -d &> /dev/null
	@sleep 5 # wait until the database has started up
	@DB_URL=$(DB_URL) PATH=${PATH}:${CURDIR}/bin go test -json ./internal/db ./internal/geocoding ./cmd/server ./cmd/dbcli

test-db: internal/db/models.go
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose up test_db -d
	sleep 5 # wait until the database has started up
	cd internal/db && DB_URL=$(DB_URL) go test || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose down test_db
	docker system prune -af &> /dev/null

test-geocoding:
	cd internal/geocoding && go test

test-server:
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose up test_db -d
	sleep 5 # wait until the database has started up
	DB_URL=$(DB_URL) go test -run ^TestHandlers || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose down test_db
	docker system prune -af &> /dev/null

test-cli: bin/dbcli
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose up test_db -d
	sleep 5 # wait until the database has started up
	cd internal/dbcli && PATH=${PATH}:${CURDIR}/bin DB_URL=$(DB_URL) go test || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose down test_db
	docker system prune -af &> /dev/null

dev-server: 
	DEV_SERVER_PORT=$(DEV_SERVER_PORT) docker compose up dev_server	

bin/server: main.go handlers.go $(wildcard internal/**/*)
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/server ./cmd/server

bin/healthcheck: cmd/healthcheck/healthcheck.go
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/healthcheck ./cmd/healthcheck

bin/dbcli: cmd/dbcli/cli.go
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/dbcli ./cmd/dbcli

migrations: bin/dbcli
	PATH=${PATH}:${CURDIR}/bin ./scripts/run-migrations.sh