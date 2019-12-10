#!/usr/bin/make -f

export GO111MODULE = on
export GOPROXY = https://goproxy.io

TIME_NOW=$$(date +"%Y-%m-%d %H:%M.%S")

.PHONY: go.sum
go.sum: go.mod
	@go mod tidy
	@go mod verify
	@go mod download

.PHONY: push
push: go.sum
	@git commit -am "UPDATE $(TIME_NOW)" && git push

.PHONY: fabric-network
fabric-network:
	@cd scripts/first-network && ./byfn.sh up

.PHONY: fabric-gen
fabric-gen:
	@cd $HOME/go/src/github.com/hyperledger/fabric && \
	make cryptogen && make configtxgen
