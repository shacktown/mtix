#Terminal 3
set -x
CC=mtix
CHAN=myc
#peer chaincode install -p chaincodedev/chaincode/mtix -n mycc -v 0
peer chaincode install -p chaincodedev/chaincode/$CC -n $CC -v 0
peer chaincode instantiate -n $CC -v 0 -c '{"Args":["",""]}' -C $CHAN
#peer chaincode instantiate -n mtix -v 0 -c '{"Args":["",""]}' -C mychannel
