BINARY=email-newsletter

.PHONY: fmt init_db run test build clean tidy

all: test

fmt:
	go fmt ./...

init_db:
	./scripts/init_db.sh

run: fmt
	go run ./cmd/.

test: fmt
	go test ./... -- -shuffle

build: fmt
	go build -o target/$(BINARY) .

clean:
	rm -rf target

tidy:
	go mod tidy
