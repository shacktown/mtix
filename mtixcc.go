package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage movie ticket sales
type SimpleAsset struct {
}

// MovieShowing represents a movie showtime including:
// showId (key), theater name, hall name,  movie name, date/time of the showing,
// 	ticket price, numSeats (max tickets), tixSold (number of tickets sold)
type MovieShowing struct {
	theater  string
	movie    string
	dateTime string
	price    string
	numSeats string
	tixSold  string
}

// Movie represents a schedule movie showing along with it's unique showID key
type Movie struct {
	showID string
	show   MovieShowing
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
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
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
	fmt.Println("show to showJSON created: ", showJSON)

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
//including the theater, hall, date, time, price, etc.
func CreateMovieShowing(args []string) (MovieShowing, error) {

	var aShow MovieShowing

	fmt.Println("CreateMovieShowing() args:  ")
	for i := 0; i < len(args); i++ {
		fmt.Println("args:", i, args[i])
	}
	/*
		if len(args) != 7 {
			fmt.Println("CreateMovieShowing(): Incorrect number of arguments. Expecting 7 ")
			return aShow, errors.New("CreateMovieShowing() : Incorrect number of arguments. Expecting 7 ")
		}*/

	aShow = MovieShowing{args[0], args[1], args[2], args[3], args[4], args[5]}
	fmt.Println("CreateMovieShowing() : Show Object : ", aShow)

	return aShow, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
