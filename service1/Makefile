LOCAL_BIN=$(CURDIR)/bin
LINTER_CMD=$(LOCAL_BIN)/golangci-lint

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
lint:
	$(LINTER_CMD) run --config=.golangci.yml ./...
lint-fix:
	$(LINTER_CMD) run --config=.golangci.yml --fix ./...

migrate:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down
deps:
	go mod tidy
	go mod download
	go mod vendor
