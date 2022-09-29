#!/bin/bash

# This script is used to compile releases

# osx
go build -o bin/packetloss .

# windows
GOOS=windows GOARCH=amd64 go build -o bin/packetloss-win64.exe .
