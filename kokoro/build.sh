#!/bin/bash
set -ex

cd git
export GOPATH=$(pwd)
echo "GOPATH=$GOPATH"
mkdir -p ${GOPATH}/bin
mkdir -p ${GOPATH}/lib
mkdir -p ${GOPATH}/src

go version

mv ./apidVerifyApiKey ./src
cd ./src/apidVerifyApiKey


echo "Getting glide"
go get github.com/Masterminds/glide
echo "Install dependencies for tests"
time ${GOPATH}/bin/glide up -v

go test $(glide novendor)
