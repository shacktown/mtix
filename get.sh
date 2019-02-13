# get the provided key from the ledger
set -x
CC=mtix
CHAN=mychannel

echo Getting information for key $1
echo ------------------------------
peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["get", "'$1'"]}'
