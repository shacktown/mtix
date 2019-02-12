# get the provided key from the ledger
set -x
CC=mtix
CHAN=myc
theater="Regal1"

echo Getting information for key $1
echo ------------------------------
peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["get", "'$1'"]}'
