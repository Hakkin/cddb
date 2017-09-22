@echo off

SET GOOS=windows
SET GOARCH=amd64
SET ldflags=-s -w
SET output=./build/windows
SET binary="cddb.exe"

go build -v -o "%output%/%binary%" -ldflags "%ldflags%" standalone.go