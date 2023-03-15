#!/bin/bash
export accessport=10000
export secret=somesecret
export JUMP_TOKEN_LIFETIME=3600
export JUMP_TOKEN_ROLE=client
export JUMP_TOKEN_SECRET=somesecret
export JUMP_TOKEN_TOPIC=123
export JUMP_TOKEN_CONNECTION_TYPE=shell
export JUMP_TOKEN_AUDIENCE=http://[::]:${accessport}
export client_token=$(../../cmd/jump/jump token)
echo "client_token=${client_token}"
export JUMP_CLIENT_LOCAL_PORT=2222
export JUMP_CLIENT_RELAY_SESSION=http://[::]:${accessport}/shell/123
export JUMP_CLIENT_TOKEN=$client_token
export JUMP_CLIENT_DEVELOPMENT=true
../../cmd/jump/jump client

