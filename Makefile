.PHONY: deps
deps:
	go mod tidy -v

.PHONY: migrate
migrate:
	goose postgres "user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

.PHONY: build-binary
build-binary:
	go build -o ./app ./cmd

.PHONY: build
build:
	deps build-binary

