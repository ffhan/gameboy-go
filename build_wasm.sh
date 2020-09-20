#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./cmd/main.wasm ./cmd/main.go
