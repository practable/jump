#!/bin/bash
export accessport=10000
export relayport=10001
export JUMP_RELAY_AUDIENCE=http://[::]:$accessport
export JUMP_RELAY_BUFFER_SIZE=128
export JUMP_RELAY_LOG_LEVEL=trace
export JUMP_RELAY_LOG_FORMAT=text
export JUMP_RELAY_LOG_FILE=/var/log/jump/jump.log
export JUMP_RELAY_PORT_ACCESS=$accessport
export JUMP_RELAY_PORT_RELAY=$relayport
export JUMP_RELAY_PROFILE=false
export JUMP_RELAY_SECRET=somesecret
export JUMP_RELAY_STATS_EVERY=5s
export JUMP_RELAY_URL=ws://[::]:${relayport}
../../cmd/jump/jump relay
