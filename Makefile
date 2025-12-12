APP_NAME := aiu
BUILD_DIR := build
BIN_DIR := $(GOPATH)/bin
MAIN_PATH := ./cmd/aiu

.PHONY: all build install clean run test fmt vet

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# Setup development environment
dev-setup:
	mkdir -p ~/.aiu/prompts
	@if [ ! -f ~/.aiu/config.yaml ]; then \
		cp config.yaml.example ~/.aiu/config.yaml; \
		echo "Created config at ~/.aiu/config.yaml"; \
	fi
	@if [ -d prompts ] && [ -n "$$(ls -A prompts)" ]; then \
		cp prompts/*.md ~/.aiu/prompts/; \
		echo "Copied sample prompts to ~/.aiu/prompts/"; \
	fi
