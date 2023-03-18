#!/bin/bash
export CWD=$(pwd)
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

function freeports(){
 #https://unix.stackexchange.com/questions/55913/whats-the-easiest-way-to-find-an-unused-local-port
	ports=$(comm -23 <(seq 49152 65535 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | sort -n | head -n 3)
}
freeports
portarray=($ports) #convert to string array
port_access=${portarray[0]}
port_relay=${portarray[1]}
port_local=${portarray[2]}

export JUMP_RELAY_ACCESS_BASE_PATH=/api/v1
export JUMP_RELAY_AUDIENCE=http://[::]:$port_access
export JUMP_RELAY_BUFFER_SIZE=128
export JUMP_RELAY_LOG_LEVEL=trace
export JUMP_RELAY_LOG_FORMAT=text
export JUMP_RELAY_LOG_FILE=stdout
export JUMP_RELAY_PORT_ACCESS=$port_access
export JUMP_RELAY_PORT_RELAY=$port_relay
export JUMP_RELAY_PROFILE=false
export JUMP_RELAY_SECRET=somesecret
export JUMP_RELAY_STATS_EVERY=5s
export JUMP_RELAY_URL=ws://[::]:${port_relay}
(cd $SCRIPT_DIR ; ../../cmd/jump/jump relay &> $CWD/relay.log &)
rpid=$1
echo "started relay"
sleep 1

#make tokens
export JUMP_TOKEN_LIFETIME=3600
export JUMP_TOKEN_ROLE=host
export JUMP_TOKEN_SECRET=$JUMP_RELAY_SECRET
export JUMP_TOKEN_TOPIC=123
export JUMP_TOKEN_CONNECTION_TYPE=connect
export JUMP_TOKEN_AUDIENCE=$JUMP_RELAY_AUDIENCE
export JUMP_HOST_TOKEN=$(cd $SCRIPT_DIR; ../../cmd/jump/jump token)
export JUMP_TOKEN_ROLE=client
export JUMP_CLIENT_TOKEN=$(cd $SCRIPT_DIR; ../../cmd/jump/jump token)

#echo $JUMP_HOST_TOKEN | decode-jwt
#echo $JUMP_CLIENT_TOKEN | decode-jwt

export JUMP_HOST_LOCAL_PORT=22
export JUMP_HOST_ACCESS=${JUMP_RELAY_AUDIENCE}${JUMP_RELAY_ACCESS_BASE_PATH}/${JUMP_TOKEN_CONNECTION_TYPE}/${JUMP_TOKEN_TOPIC}
export JUMP_HOST_DEVELOPMENT=true
(cd $SCRIPT_DIR ; ../../cmd/jump/jump host &> $CWD/host.log &)
hpid=$1
echo "started host"
sleep 1

export JUMP_CLIENT_LOCAL_PORT=$port_local
export JUMP_CLIENT_RELAY_SESSION=${JUMP_RELAY_AUDIENCE}${JUMP_RELAY_ACCESS_BASE_PATH}/${JUMP_TOKEN_CONNECTION_TYPE}/${JUMP_TOKEN_TOPIC}
export JUMP_CLIENT_DEVELOPMENT=true
(cd $SCRIPT_DIR ; ../../cmd/jump/jump client &> $CWD/client.log &)
cpid=$1
echo "started client"
sleep 1

ssh -p $port_local $USER@127.0.0.1
kill $rpid
kill $cpid
kill $hpid
