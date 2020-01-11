#!/usr/bin/make -f

TIME_NOW=$$(date +"%Y-%m-%d %H:%M.%S")
PROXY=GOPROXY=https://goproxy.io

.PHONY: build
build:
	go build -mod=readonly -o build/fabriclid ./cmd/fabriclid

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

.PHONY: go.sum
go.sum: go.mod
	@$(PROXY) go mod tidy
	@$(PROXY) go mod download
	@$(PROXY) go mod verify

.PHONY: install
install:
	@go install -v ./cmd/fabriclid

.PHONY: fabric-network
fabric-network:
	@cd scripts/first-network && ./byfn.sh up

.PHONY: fabric-gen
fabric-gen:
	@cd $HOME/go/src/github.com/hyperledger/fabric && \
	make cryptogen && make configtxgen

.PHONY: push
push: go.sum
	@git commit -am "UPDATE $(TIME_NOW)" && git push
