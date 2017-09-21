#!/bin/bash
set -ex

BUILDROOT=${BUILDROOT:-github/apidVerifyApiKey}
export BUILDROOT

# Make a temporary GOPATH to build in
gobase=`mktemp -d`
base=${gobase}/src/github.com/apid/apidVerifyApiKey
GOPATH=${gobase}
export GOPATH

base=${GOPATH}/src/github.com/apid/apidVerifyApiKey
mkdir -p ${base}
(cd ${BUILDROOT}; tar cf - .) | (cd ${base}; tar xf -)
cd ${base}


echo "Getting glide"
go get github.com/Masterminds/glide
echo "Install dependencies for tests"
time ${GOPATH}/bin/glide up -v

go version

go test $(glide novendor)
