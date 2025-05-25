BINARY=email-newsletter

.PHONY: fmt init_db run test test-race build clean tidy

all: test

fmt:
	go fmt ./...

init_db:
	./scripts/init_db.sh

run: fmt
	@LOG_LEVEL=trace go run .

test: fmt
	@GIN_MODE=release go test ./... -- -shuffle

test-race: fmt
	@go test -race ./...

build: fmt
	@echo Build email-newsletter
	go build -tags netgo -ldflags '-s -w' -o target/$(BINARY) .

clean:
	rm -rf target

tidy:
	@go mod tidy
