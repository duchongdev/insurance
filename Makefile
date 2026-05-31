.PHONY: tidy test build run

GOPROXY ?= https://goproxy.io,direct
GOSUMDB ?= sum.golang.org

tidy:
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) go mod tidy

test: tidy
	go test ./...

build: tidy
	go build -o bin/bridge ./cmd/server

run: build
	CONFIG_PATH=config/config.yaml ./bin/bridge
