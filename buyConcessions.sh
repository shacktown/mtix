# buy concessions water, soda or popcorn
# For a specific each movie time (3:00, 5:00, 7:00, 9:00)
#   for this example the theater is Regal1
set -x
CC=mtix
CHAN=myc
theater="Regal1"
peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["ConcessionsAvailable", "Regal1", "soda", "4", "2019-02-14T21:00:00-05:00"]}'
peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["BuyConcession", "Regal1", "soda", "4", "2019-02-14T21:00:00-05:00"]}'
#peer chaincode invoke  -n mycc -C mychannel -c '{"Args":["get", "buyTxn1"]}'