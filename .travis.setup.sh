#!/usr/bin/env sh

mkdir -p "$GOPATH/src/github.com/Boostport"
ln -s /source $GOPATH/src/github.com/Boostport/avatica
apk update
apk add git
go get -u github.com/kardianos/govendor
go get -u github.com/go-playground/overalls
go get -u github.com/mattn/goveralls
