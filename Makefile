
build: test
	go build -o ./bin/app ./cmd

dev: build
	./cmd/dev.sh

.PHONY: test
test:
	go test -v ./...