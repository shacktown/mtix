#Terminal 2
set -x
#cd mtix
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mtix:0 ./mtix
#CORE_PEER_ADDRESS=peer0.org1.example.com:7052 CORE_CHAINCODE_ID_NAME=mtix:1 ./mtix