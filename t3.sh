#Terminal 3
set -x
peer chaincode install -p chaincodedev/chaincode/mtix -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":["",""]}' -C myc
