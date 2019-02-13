# buy movie tickets 
set -x
CC=mtix
CHAN=mychannel

peer chaincode invoke -n $CC -C $CHAN -c '{"Args":["BuyTix", "show1", "4", "window1"]}'
#peer chaincode invoke  -n mycc -C mychannel -c '{"Args":["get", "buyTxn1"]}'