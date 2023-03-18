#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
(cd $SCRIPT_DIR; rm -rf ../../internal/access/models ; rm -rf ../../internal/access/restapi ; swagger generate server -t ../../internal/access -f ../../api/access.yml --exclude-main -A access ; )
go mod tidy

