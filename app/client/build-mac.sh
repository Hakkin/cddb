#!/bin/bash

GOOS="darwin"
GOARCH="amd64"
ldflags="-s -w"
output="./build/darwin/"
binary="cddb"

export GOOS="$GOOS"
export GOARCH="$GOARCH"

go build -v -o "$output/$binary" -ldflags "$ldflags" standalone.go