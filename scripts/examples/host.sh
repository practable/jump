#!/bin/bash
export accessport=10000
export relayport=10001
export JUMP_TOKEN_LIFETIME=3600
export JUMP_TOKEN_ROLE=host
export JUMP_TOKEN_SECRET=somesecret
export JUMP_TOKEN_TOPIC=123
export JUMP_TOKEN_CONNECTION_TYPE=shell
export JUMP_TOKEN_AUDIENCE=http://[::]:${accessport}
export host_token=$(../../cmd/jump/jump token)
echo "host_token=${host_token}"
export JUMP_HOST_LOCAL_PORT=22
export JUMP_HOST_ACCESS=http://[::]:${accessport}/shell/123
export JUMP_HOST_TOKEN=$host_token
export JUMP_HOST_DEVELOPMENT=true
../../cmd/jump/jump host 
