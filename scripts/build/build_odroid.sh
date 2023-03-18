#!/bin/bash
export GOOS=linux
export GOARCH=arm64
now=$(date +'%Y-%m-%d_%T')
(cd ../../cmd/jump; go build -ldflags "-X 'github.com/practable/jump/cmd/jump/cmd.Version=`git describe --tags`' -X 'github.com/practable/jump/cmd/jump/cmd.BuildTime=$now'"; ./jump version)
