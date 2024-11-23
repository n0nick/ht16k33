#!/usr/bin/env sh

GOOS=linux GOARCH=arm CC_FOR_TARGET=arm-linux-gnueabi-gcc go build ./main.go
