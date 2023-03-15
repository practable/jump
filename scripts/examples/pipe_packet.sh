#!/bin/bash
export JUMP_PIPE_PACKET_PORT_LISTEN=2222
export JUMP_PIPE_PACKET_PORT_TARGET=22
export JUMP_PIPE_PACKET_LOG_LEVEL=trace
export JUMP_PIPE_PACKET_LOG_FORMAT=json
../../cmd/jump/jump pipe packet
