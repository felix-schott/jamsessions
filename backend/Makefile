include .env

.PHONY = test-ci test-db test-geocoding test-server test-cli

TEST_DB_PORT = 54930
DEV_SERVER_PORT = 36230

DB_URL := "host=localhost port=$(TEST_DB_PORT) user=read_write password=$(READ_WRITE_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable"

internal/db/models.go: internal/db/schema.sql internal/db/query.sql internal/db/sqlc.yaml 
	sqlc generate -f internal/db/sqlc.yaml
	# for some reason sqlc isn't working correctly here, below is a hacky fix
	sed -i 's/d\.d/d.dates/g' internal/db/query.sql.go

test-ci: internal/db/models.go bin/dbcli # run in ci pipeline
	@TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml up test_db -d
	@sleep 7 # wait until the database has started up
	@DB_URL=$(DB_URL) PATH=${PATH}:${CURDIR}/bin go test ./internal/db ./internal/geocoding ./cmd/server ./cmd/dbcli

test-db: internal/db/models.go
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml up test_db -d
	sleep 7 # wait until the database has started up
	cd internal/db && DB_URL=$(DB_URL) go test || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml down test_db
	docker system prune -af &> /dev/null

test-geocoding:
	cd internal/geocoding && go test

test-server:
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml up test_db -d
	sleep 5 # wait until the database has started up
	cd cmd/server && DB_URL=$(DB_URL) go test -run ^TestHandlers || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml down test_db
	docker system prune -af &> /dev/null

test-cli: bin/dbcli
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml up test_db -d
	sleep 5 # wait until the database has started up
	cd cmd/dbcli && PATH=${PATH}:${CURDIR}/bin DB_URL=$(DB_URL) go test || true
	TEST_DB_PORT=$(TEST_DB_PORT) docker compose -f test.docker-compose.yml down test_db
	docker system prune -af &> /dev/null

bin/server: cmd/server/main.go cmd/server/handlers.go $(wildcard internal/**/*)
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/server ./cmd/server

bin/healthcheck: cmd/healthcheck/healthcheck.go
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/healthcheck ./cmd/healthcheck

bin/dbcli: cmd/dbcli/cli.go
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/dbcli ./cmd/dbcli