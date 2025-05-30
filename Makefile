PROTO_DIR := api/proto
OUT_DIR := pkg/grpc

PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: generate

generate: $(PROTO_FILES)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILES) \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \

.PHONY: run

run:
	go run ./cmd/notification/

.PHONY: migrate

migrate:
	go run ./cmd/migrate/

.PHONY: compose-up
compose-up:
	docker compose -f docker-compose.yml up 

.PHONY: compose-build
compose-build:
	docker compose -f docker-compose.yml build