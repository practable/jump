#!/bin/bash
export JUMP_PIPE_DEVELOP_PORT_LISTEN=2222
export JUMP_PIPE_DEVELOP_PORT_TARGET=22
export JUMP_PIPE_DEVELOP_LOG_LEVEL=trace
export JUMP_PIPE_DEVELOP_LOG_FORMAT=json
../../cmd/jump/jump pipe develop
