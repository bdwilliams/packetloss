#!/bin/bash

# This script is used to compile releases

# osx
go build -o bin/packetloss .

# windows
GOOS=windows GOARCH=amd64 go build -o bin/packetloss-win64.exe .

# linux (amd64)
GOOS=linux GOARCH=amd64 go build -o bin/packetloss-amd64-linux .

# linux (arm64)
GOOS=linux GOARCH=arm64 go build -o bin/packetloss-arm64-linux .
