.PHONY: tidy test build run package install start stop deploy-test admin-build admin-dev

GOPROXY ?= https://goproxy.io,direct
GOSUMDB ?= sum.golang.org

admin-build:
	bash scripts/build-admin.sh

admin-dev:
	cd web/admin && npm run dev

tidy:
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) go mod tidy

test: tidy
	go test ./...

build: tidy
	go build -o bin/bridge ./cmd/server

run: build
	CONFIG_PATH=config/config.yaml ./bin/bridge

# 生成交付压缩包 → dist/insurance-bridge-<go版本>-<日期>.tar.gz
package:
	bash scripts/package.sh

install:
	bash scripts/install.sh

start:
	bash scripts/start.sh

stop:
	bash scripts/stop.sh

deploy-test:
	bash scripts/deploy-test.sh
