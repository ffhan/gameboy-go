#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./cmd/wasm/main.wasm ./cmd/wasm/main.go
