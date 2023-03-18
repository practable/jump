#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
( cd $SCRIPT_DIR; rm ../../internal/access/restapi/configure_access.go)
