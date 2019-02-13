# stock concessions for a theater 
#   for this example the theater is Regal1
set -x
CC=mtix
CHAN=mychannel
theater="Regal1"

peer chaincode invoke  -n $CC -C $CHAN -c '{"Args":["StockConcession", "'$theater'", "soda",  "100", "7.00"]}'
peer chaincode invoke  -n $CC -C $CHAN -c '{"Args":["StockConcession", "'$theater'", "water",  "100", "5.00"]}'
peer chaincode invoke  -n $CC -C $CHAN -c '{"Args":["StockConcession", "'$theater'", "popcorn",  "100", "8.00"]}'