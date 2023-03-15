#!/bin/sh

# Make token
export accessport=10000
export JUMP_TOKEN_LIFETIME=30
export JUMP_TOKEN_ROLE=stats
export JUMP_TOKEN_SECRET=somesecret
export JUMP_TOKEN_TOPIC=stats
export JUMP_TOKEN_CONNECTION_TYPE=shell
export JUMP_TOKEN_AUDIENCE=http://[::]:$accessport
export client_token=$(../../cmd/jump/jump token)
# echo "client_token=${client_token}"

# Request Access
export ACCESS_URL="${SHELLTOKEN_AUDIENCE}/shell/stats"

export STATS_URL=$(curl -X POST  \
-H "Authorization: ${client_token}" \
$ACCESS_URL | jq -r '.uri')

# echo $STATS_URL
# Connect to stats channel & issue {"cmd":"update"}

echo '{"cmd":"update"}' | websocat -B 1048576 -n1 "$STATS_URL" | jq .
