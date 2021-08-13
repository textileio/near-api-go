include .bingo/Variables.mk

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run
.PHONYY: lint

build:
	go build ./...
.PHONY: build

test:
	go test ./...
.PHONY: test
