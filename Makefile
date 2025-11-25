APP_NAME := aiu
BIN_DIR := $(GOPATH)/bin
MAIN_PATH := ./cmd/aiu

.PHONY: all build install clean run test fmt vet

all: build

build:
	go build -o $(APP_NAME) $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	rm -f $(APP_NAME)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# Setup development environment
dev-setup:
	mkdir -p ~/.config/ai-utils/prompts
	@if [ ! -f ~/.config/ai-utils/config.yaml ]; then \
		cp config.yaml.example ~/.config/ai-utils/config.yaml; \
		echo "Created config at ~/.config/ai-utils/config.yaml"; \
	fi
	@if [ -d prompts ] && [ -n "$$(ls -A prompts)" ]; then \
		cp prompts/*.md ~/.config/ai-utils/prompts/; \
		echo "Copied sample prompts to ~/.config/ai-utils/prompts/"; \
	fi
