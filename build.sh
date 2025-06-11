#!/bin/bash
set -e

mkdir -p bin

echo "Building selector-tester for linux/amd64..."
env GOOS=linux GOARCH=amd64 go build -o bin/selector-tester_linux_amd64 server.go

echo "Building selector-tester for darwin/amd64..."
env GOOS=darwin GOARCH=amd64 go build -o bin/selector-tester_darwin_amd64 server.go

echo "Done."
