#!/bin/bash
set -e
export GOPATH="${HOME}/go"
export PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"
export GHC_APP_GOPATH="${GOPATH}/src/ghc-app"
pushd "${GHC_APP_GOPATH}" > /dev/null
go get -u github.com/jteeuwen/go-bindata/...
go generate ghc-app # run go-bindata
go get -d ./...
export GOBIN='/srv/bin'
go install ghc-app
popd

export EVENT_DIGEST_PATH="${GOPATH}/src/event-digest"
pushd "${EVENT_DIGEST_PATH}" > /dev/null
go get -d ./...
go install event-digest
popd
