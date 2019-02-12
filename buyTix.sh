# buy movie tickets 
# For each movie time (3:00, 5:00, 7:00, 9:00)
#   schedule a movie in each hall (1-5)
#   for this example the theater is Regal1
#   an incrementing 'showID' is created as the unique key for each movie showing
set -x
CC=mtix
CHAN=myc
theater="Regal1"
ID=0
peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["BuyTix", "show1", "4", "window1"]}'
#peer chaincode invoke  -n mycc -C mychannel -c '{"Args":["get", "buyTxn1"]}'