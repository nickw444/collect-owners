#!/bin/sh
set -eux

mkdir -p release

VERSION=$(git describe --dirty --always)
LDFLAGS="-X main.Version=$VERSION"

GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o release/collect-owners-linux-amd64
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o release/collect-owners-linux-arm64
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o release/collect-owners-darwin-amd64
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o release/collect-owners-windows-amd64
