BINARY_NAME := slpx
BUILD_DIR := build
CMD_DIR := cmd/slpx
PKG_DIRS := ./pkg/...
GO := go
GOFLAGS := 
LDFLAGS := -s -w

SRCS := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGET := $(BUILD_DIR)/$(BINARY_NAME)

.DEFAULT_GOAL := all

.PHONY: all build test clean fmt vet lint run install help

all: build

build: $(TARGET)

$(TARGET): $(SRCS)
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET) ./$(CMD_DIR)

test:
	$(GO) test -v -race -cover $(PKG_DIRS)

clean:
	rm -rf $(BUILD_DIR)
	$(GO) clean

fmt:
	$(GO) fmt $(PKG_DIRS)

vet:
	$(GO) vet $(PKG_DIRS)

lint: fmt vet

run: build
	./$(TARGET)

install: build
	$(GO) install ./$(CMD_DIR)

help:
	@echo "Available targets:"
	@echo "  all      - Build the project (default)"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests with race detector and coverage"
	@echo "  clean    - Remove build artifacts"
	@echo "  fmt      - Format Go source files"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run fmt and vet"
	@echo "  run      - Build and run the binary"
	@echo "  install  - Install the binary to GOPATH/bin"
	@echo "  help     - Show this help message"

