package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage movie ticket sales
type SimpleAsset struct {
}

// MovieShowing represents a movie showtime including:
// showId (key), theater name, Hall name,  movie name, date/time of the showing,
// 	ticket price, numSeats (max tickets), tixSold (number of tickets sold)
type MovieShowing struct {
	Theater  string
	Hall     string
	Movie    string
	DateTime string
	Price    string
	NumSeats string
	TixSold  string
}

// Movie represents a schedule movie showing along with it's unique showID key
type Movie struct {
	ShowID string
	Show   MovieShowing
}

// Init is called during chaincode instantiation to initialize data
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	/*   args := stub.GetStringArgs()
	if len(args) != 2 {
	        return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	*/
	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	/*   err := stub.PutState(args[0], []byte(args[1]))
	    if err != nil {
	            return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
		}
	*/
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "ScheduleShow" {
		result, err = ScheduleMovieShowing(stub, args)
	} else if fn == "BuyTix" {
		result, err = BuyTix(stub, args)
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

/*
TicketPurchase (created with transactions):
txnId   showId  quantity    purchaseTime                    revenue             numSodas    seller
1       1       4           2019-02-05T13:00:00-05:00   quantity * price    4           W1
*/
type TicketPurchase struct {
	TxnID        string
	ShowID       string
	Quantity     string
	PurchaseTime string
	Revenue      string
	NumSodas     string
	Seller       string
}

// BuyTix purchases 1 or more tickets (quantity) for a specific movie showing
// from a specified seller (window number). This reduces the number of tix available
// The purchase transaction is added to the ledger.
// If the key exists, it will override the value with the new one
func BuyTix(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	printArgs("BuyTix", args)

	if len(args) != 3 {
		fmt.Println("BuyTix(): Incorrect number of arguments. Expecting 3")
		return "", errors.New("BuyTix() : Incorrect number of arguments. Expecting 3 ")
	}

	var buytix TicketPurchase
	//TODO get the txn ID or produce it
	buytix.TxnID = "buyTxn1"

	// args:  showId, quantity, seller
	buytix.ShowID = args[0]
	buytix.Quantity = args[1]
	buytix.Seller = args[2]

	// TODO - get dateTime
	buytix.PurchaseTime = "2019-02-05T13:00:00-05:00"
	//TODO compute revenue
	buytix.Revenue = "100.00"
	buytix.NumSodas = "0"
	if buytix.Seller == "window1" {
		buytix.NumSodas = buytix.Quantity
	}
	fmt.Println("BuyTix() - purchase attributes : ", buytix)
	buytixJSON, err := json.Marshal(buytix)
	if err != nil {
		fmt.Println("buytix to buytixJSON error: ", err)
		return "", err
	}
	fmt.Println("buytix to buytixJSON created: ", string(buytixJSON))

	//get show information, verify tickets are available and update tickets sold
	bytes, err := stub.GetState(buytix.ShowID)
	if err != nil {
		return "", fmt.Errorf("BuyTix failed to get show information for: %s", buytix.ShowID)
	}
	fmt.Println("---------------- Here is the byte string", string(bytes))

	var movieShowing MovieShowing
	err = json.Unmarshal(bytes, &movieShowing)
	var q, _ = strconv.Atoi(buytix.Quantity)
	var n, _ = strconv.Atoi(movieShowing.NumSeats)
	var s, _ = strconv.Atoi(movieShowing.TixSold)
	if q > (n - s) {
		return "", fmt.Errorf("Failed to complete ticket purchase - not enough tickets left. txnId=%s", buytix.TxnID)
	}
	//update the show with the new number of purchased tickets
	movieShowing.TixSold = fmt.Sprintf("%d", s+q)
	showJSON, err := json.Marshal(movieShowing)
	if err != nil {
		fmt.Println("movieShowing to showJSON error: ", err)
		return "", err
	}
	err = stub.PutState(buytix.ShowID, showJSON)
	if err != nil {
		return "", fmt.Errorf("BuyTix failed to update purchased quantify for: %s", buytix.ShowID)
	}
	err = stub.PutState(buytix.TxnID, buytixJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to complete ticket purchase for: %s", string(buytixJSON))
	}
	return string(buytixJSON), nil

}

// ScheduleMovieShowing schedules and stores the provided movie showing on the ledger.
// If the key exists, it will override the value with the new one
func ScheduleMovieShowing(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	show, err := CreateMovieShowing(args[0:])
	if err != nil {
		return "", err
	}
	showJSON, err := json.Marshal(show)
	if err != nil {
		fmt.Println("show to showJSON error: ", err)
		return "", err
	}
	fmt.Println("show to showJSON created: ", string(showJSON))

	err = stub.PutState("show1", showJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return "show1", nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

//CreateMovieShowing creates a new movie showing object
//including the Theater, Hall, date, time, price, etc.
func CreateMovieShowing(args []string) (MovieShowing, error) {

	var aShow MovieShowing

	printArgs("CreateMovieShowing", args)

	if len(args) != 7 {
		fmt.Println("CreateMovieShowing(): Incorrect number of arguments. Expecting 7 ")
		return aShow, errors.New("CreateMovieShowing() : Incorrect number of arguments. Expecting 7 ")
	}

	aShow = MovieShowing{args[0], args[1], args[2], args[3], args[4], args[5], args[6]}
	fmt.Println("CreateMovieShowing() - Show attributes : ", aShow)

	return aShow, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

func printArgs(caller string, args []string) {
	fmt.Println(caller, "() args:  ")
	for i := 0; i < len(args); i++ {
		fmt.Println("args:", i, args[i])
	}
}
