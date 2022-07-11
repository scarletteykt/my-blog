.PHONY: deps
deps:
	go mod tidy -v

.PHONY: migrate
migrate:
	goose --dir /app/migration postgres "postgres://$(DB_USER):$(DB_PASSWORD)@database:5432/$(DB_NAME)?sslmode=disable" status
	goose --dir /app/migration postgres "postgres://$(DB_USER):$(DB_PASSWORD)@database:5432/$(DB_NAME)?sslmode=disable" up

.PHONY: build-binary
build-binary:
	go build -o /app/$(OUTPUT_BINARY) /app/$(BUILD_DIR)

.PHONY: build
build: deps build-binary

.PHONY: run
run:
	/app/$(OUTPUT_BINARY)