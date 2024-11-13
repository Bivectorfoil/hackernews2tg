#!/bin/bash

export CGO_ENABLED=0

mkdir -p dist

LDFLAGS="-w -s"

echo "Building linux version..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "dist/newsboy"

echo "Build completed!"
echo "Linux version: dist/newsboy"

chmod +x dist/newsboy
