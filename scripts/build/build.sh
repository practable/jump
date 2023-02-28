#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation issues in go build command
export GOOS=linux
now=$(date +'%Y-%m-%d_%T')
(cd ../../cmd/jump; go build -ldflags "-X 'github.com/practable/jump/cmd/jump/cmd.Version=`git describe --tags`' -X 'github.com/practable/jump/cmd/jump/cmd.BuildTime=$now'"; ./jump version)

