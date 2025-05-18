BINARY=email-newsletter

.PHONY: fmt init_db run test test-race build clean tidy

all: test

fmt:
	go fmt ./...

init_db:
	./scripts/init_db.sh

run: fmt
	LOG_LEVEL=trace go run .

test: fmt
	go test ./... -- -shuffle

test-race: fmt
	go test -race ./...

build: fmt
	go build -o target/$(BINARY) .

clean:
	rm -rf target

tidy:
	go mod tidy
