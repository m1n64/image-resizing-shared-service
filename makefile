APP_NAME := imageresizer
CMD_PATH := ./cmd/server
BIN_PATH := ./tmp/$(APP_NAME)

.PHONY: help build run prod dev debug proto

help:
	@echo ""
	@echo "🔥  Available Makefile Commands 🔥"
	@echo "--------------------------------------------------------------"
	@echo "💻  Start environment:"
	@echo "  make build          - Build project"
	@echo "  make run            - Run project"
	@echo "  make prod           - Build production project"
	@echo "  make dev            - Run project by Air (required Air and Delve)"
	@echo "  make debug          - Run project by Delve"
	@echo "  make clean          - Clean tmp dir"
	@echo "  make proto          - Compile protobuf (require protoc and protoc-gen-go)"
	@echo ""

build:
	CGO_ENABLED=1 go build -o $(BIN_PATH) $(CMD_PATH)

run:
	go run $(CMD_PATH)/main.go

prod:
	CGO_ENABLED=1 go build -ldflags "-X main.GinMode=release" -o $(BIN_PATH) $(CMD_PATH)

dev:
	air

debug:
	dlv debug $(CMD_PATH)/main.go --headless --listen=:2345 --api-version=2

clean:
	rm -rf ./tmp

proto:
	protoc --go_out=. --go-grpc_out=. proto/image_service.proto