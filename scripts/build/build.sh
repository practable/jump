#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation issues in go build command
export GOOS=linux
now=$(date +'%Y-%m-%d_%T')
go build -ldflags "-X 'github.com/timdrysdale/gradex-cli/cmd.Version=`git describe`' -X 'github.com/timdrysdale/gradex-cli/cmd.BuildTime=$now'"  

