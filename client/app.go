package main

/*
	- CORE_PEER_ID=example02
 	- CORE_PEER_ADDRESS=peer:7051
	- CORE_PEER_LOCALMSPID=DEFAULT
	orderer:7050     127.0.0.1:7050
	   - CORE_PEER_ID=example02
	  - CORE_PEER_ADDRESS=peer:7051

	  CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0

	  organization Name: SampleOrg

*/
import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel" //"devopschannel"
	orgName        = "org1"      //"SampleOrg"
	orgAdmin       = "Admin"
	ordererOrgName = "Orderer"
	ccID           = "mycc" //"devopschannel-example_cc2"
)

func main() {
	configPath := "./config1.yaml"
	configOpt := config.FromFile(configPath)

	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		fmt.Println("Failed to create new SDK: \n", err)
	}
	defer sdk.Close()

	org1ChannelClientContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	channelClient, err := channel.New(org1ChannelClientContext)
	if err != nil {
		fmt.Printf("Failed to create new channel client: %s\n", err)
	}

	var args = [][]byte{[]byte("get"),
		[]byte("show1"),
	}
	res, err := channelClient.Query(channel.Request{
		ChaincodeID: ccID,
		Fcn:         "invoke",
		Args:        args,
	})

	if err != nil {
		fmt.Printf("Failed to query: %s\n", err)

	}
	fmt.Println(string(res.Payload))

	// eventID := ".*"

	// // // Register chaincode event (pass in channel which receives event details when the event is complete)
	// reg, notifier, err := channelClient.RegisterChaincodeEvent(ccID, eventID)
	// if err != nil {
	//  fmt.Printf("Failed to register cc event: %s", err)
	// }
	//  defer channelClient.UnregisterChaincodeEvent(reg)
	/*
		res, err = channelClient.Execute(channel.Request{
			ChaincodeID: ccID,
			Fcn:         "invoke",
			Args: [][]byte{
				[]byte("move"),
				[]byte("a"),
				[]byte("b"),
				[]byte("100"),
			},
		})

		if err != nil {
			fmt.Printf("Failed to invoke: %s\n", err)

		}
		fmt.Println(string(res.Payload))
	*/
	// select {
	// case ccEvent := <-notifier:
	//  log.Printf("Received CC event: %#v\n", ccEvent)
	// case <-time.After(time.Second * 20):
	//  log.Printf("Did NOT receive CC event for eventId(%s)\n", eventID)
	// }

}
