.PHONY: build test

build:
	go build ./...

test: bin/sqlc-gen-c.wasm
	go test ./...
	
test-update: bin/sqlc-gen-c.wasm
	go test ./internal/endtoend -update

all: bin/sqlc-gen-c bin/sqlc-gen-c.wasm

bin/sqlc-gen-c: bin go.mod go.sum $(wildcard **/*.go)
	cd plugin && go build -o ../bin/sqlc-gen-c ./main.go

bin/sqlc-gen-c.wasm: bin/sqlc-gen-c
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-gen-c.wasm main.go

bin:
	mkdir -p bin
