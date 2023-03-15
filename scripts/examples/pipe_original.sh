#!/bin/bash
export JUMP_PIPE_ORIGINAL_PORT_LISTEN=2222
export JUMP_PIPE_ORIGINAL_PORT_TARGET=22
export JUMP_PIPE_ORIGINAL_LOG_LEVEL=debug
export JUMP_PIPE_ORIGINAL_LOG_FORMAT=json
../../cmd/jump/jump pipe original
