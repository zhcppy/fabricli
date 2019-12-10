#!/usr/bin/env bash

# shell is funny

set -e

IMAGE_NAME=("ca" "javaenv" "tools" "ccenv" "orderer" "peer"
"zookeeper" "kafka" "couchdb" "baseimage" "baseos")

export http_proxy=socks5://127.0.0.1:1086
export https_proxy=socks5://127.0.0.1:1086

for name in ${IMAGE_NAME[@]} ; do
    echo "docker pull hyperledger/fabric-${name}:latest"
    docker pull "hyperledger/fabric-${name}:latest"
done
