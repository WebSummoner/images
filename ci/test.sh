#!/bin/bash

export GO111MODULE="on"
go test -race -v -coverprofile=coverage.out -covermode=atomic ./...
cp coverage.out coverage.txt

go install golang.org/x/vuln/cmd/govulncheck@latest
"$(go env GOPATH)"/bin/govulncheck ./...
