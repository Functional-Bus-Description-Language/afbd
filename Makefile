NAME=afbd

.PHONY: default all build help fmt vet install uninstall

default: build

help:
	@echo "Build targets:"
	@echo "  all      Run fmt vet build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality targets:"
	@echo "  fmt   Format files with go fmt."
	@echo "  lint  Lint files with golangci-lint."
	@echo "Test targets:"
	@echo "  test  Run go test."
	@echo "  tb    Run testbenches."
	@echo "Other targets:"
	@echo "  help  Print help message."
	@echo "  go-update-deps "
	@echo "       Update go dependencies."

# Build targets
all: lint fmt build

build:
	go build -v -o $(NAME) ./cmd/$(NAME)


# Quality targets
fmt:
	go fmt ./...

lint:
	golangci-lint run


# Test targets
test:
	go test ./...

tb:
	hbs test afbd


# Installation targets
install:
	cp $(NAME) /usr/bin

uninstall:
	rm /usr/bin/$(NAME)

# Other targets:
go-update-deps:
	go get -u ./...
