package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
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
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	result := "Success"
	var err error
	if fn == "ScheduleShow" {
		result, err = ScheduleMovieShowing(stub, args)
	} else if fn == "BuyTix" {
		result, err = BuyTix(stub, args)
	} else if fn == "StockConcession" {
		result, err = StockConcession(stub, args)
	} else if fn == "BuyConcession" {
		result, err = BuyConcession(stub, args)
	} else if fn == "TicketsAvailable" {
		_, err = TicketsAvailable(stub, args)
	} else if fn == "ConcessionsAvailable" {
		_, err = ConcessionsAvailable(stub, args)
	} else if fn == "SodasAvailable" {
		_, err = SodasAvailable(stub, args)
	} else if fn == "ExchangeWaterSoda" {
		result, err = ExchangeWaterSoda(stub, args)
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

	buytix.TxnID = stub.GetTxID()

	// args:  showId, quantity, seller
	buytix.ShowID = args[0]
	buytix.Quantity = args[1]

	//first, verify that enough tickets are available
	showing, err := TicketsAvailable(stub, args[0:2])
	if err != nil {
		//not enough tickets available
		fmt.Println(err)
		return "", err
	}

	buytix.Seller = args[2]

	// Set current date / time as purchase date
	t := time.Now()
	nowFormatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	buytix.PurchaseTime = nowFormatted //"2019-02-05T13:00:00-05:00"
	//Compute revenue for this sale
	var q, _ = strconv.Atoi(buytix.Quantity)
	var p, _ = strconv.ParseFloat(showing.Price, 64)
	var s, _ = strconv.Atoi(showing.TixSold)
	rev := float64(q) * p
	buytix.Revenue = fmt.Sprintf("%f", rev)
	fmt.Println("BuyTix() - purchase attributes : ", buytix)
	buytixJSON, err := json.Marshal(buytix)
	if err != nil {
		fmt.Println("buytix to buytixJSON error: ", err)
		return "", err
	}
	fmt.Println("buytix to buytixJSON created: ", string(buytixJSON))

	//update the show with the new number of purchased tickets
	showing.TixSold = fmt.Sprintf("%d", s+q)
	showJSON, err := json.Marshal(showing)
	if err != nil {
		fmt.Println("movieShowing to showJSON error: ", err)
		return "", err
	}
	err = stub.PutState(buytix.ShowID, showJSON)
	if err != nil {
		return "", fmt.Errorf("BuyTix failed to update purchased quantity for: %s", buytix.ShowID)
	}
	err = stub.PutState(buytix.TxnID, buytixJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to complete ticket purchase for: %s", string(buytixJSON))
	}
	results := fmt.Sprintf("Buy Tickets created ID=%s Movie=%s Quantity=%s Tickets sold=%s", buytix.TxnID, showing.Movie, buytix.Quantity, showing.TixSold)
	return results, nil
}

// ScheduleMovieShowing schedules and stores the provided movie showing on the ledger.
// If the key exists, it will override the value with the new one
func ScheduleMovieShowing(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	txID := args[0]
	show, err := CreateMovieShowing(args[1:])
	if err != nil {
		return "", err
	}
	showJSON, err := json.Marshal(show)
	if err != nil {
		fmt.Println("show to showJSON error: ", err)
		return "", err
	}
	fmt.Println("show to showJSON created: ", string(showJSON))

	err = stub.PutState(txID, showJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	results := fmt.Sprintf("Show created ID=%s Movie=%s", txID, show.Movie)
	return results, nil
}

// TicketsAvailable verifies a quantity of tickets are available for a given show
func TicketsAvailable(stub shim.ChaincodeStubInterface, args []string) (MovieShowing, error) {

	var show MovieShowing
	if len(args) != 2 {
		fmt.Println("TicketsAvailable(): Incorrect number of arguments. Expecting 3, received", len(args))
		return show, fmt.Errorf("TicketsAvailable(): Incorrect number of arguments. Expecting 3, received %d", len(args))
	}
	// args:  showId, quantity
	showID := args[0]
	quantity := args[1]

	bytes, err := stub.GetState(showID)
	if err != nil {
		return show, fmt.Errorf("TicketsAvailable failed to get show information for: %s", showID)
	}

	err = json.Unmarshal(bytes, &show)
	var q, _ = strconv.Atoi(quantity)
	var n, _ = strconv.Atoi(show.NumSeats)
	var s, _ = strconv.Atoi(show.TixSold)
	var tixAvailable = n - s

	fmt.Println("TicketsAvailable requesting ", quantity, " number available is ", tixAvailable)
	if q > tixAvailable {
		return show, fmt.Errorf("Error - Out of tickets: Attempting to buy %s tickets for %s, but only %d are available", quantity, showID, tixAvailable)
	}
	return show, nil
}

//CreateMovieShowing creates a new movie showing object
//including the Theater, Hall, date, time, price, etc.
func CreateMovieShowing(args []string) (MovieShowing, error) {

	var aShow MovieShowing

	printArgs("CreateMovieShowing", args)

	if len(args) != 7 {
		msg := fmt.Sprintf("CreateMovieShowing(): Incorrect number of arguments. Expecting 7 received %d", len(args))
		return aShow, errors.New(msg)
	}

	aShow = MovieShowing{args[0], args[1], args[2], args[3], args[4], args[5], args[6]}
	fmt.Println("CreateMovieShowing() - Show attributes : ", aShow)

	return aShow, nil
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

// Concession represents a refreshment for sale including:
// ConID (key), soda, water or popcorn
// Inventory, the # of units in stock
// Price, sales prices
type Concession struct {
	Theater   string
	ConID     string
	Inventory string
	Price     string
}

// ConPurchase represents the purchase of a concession
// TxnID - unique key
// ConID         type of concession
// Quantity     # of units purchased
// PurchaseTime day/time of txn
// Revenue      total revenue for sale
// ShowTime     the movie show time
type ConPurchase struct {
	TxnID        string
	Theater      string
	ConID        string
	Quantity     string
	Revenue      string
	PurchaseTime string
}

// SodaCount represents the number of Soda's purchased for each showTime
//	a limit of 200 sodas per show is enforced
type SodaCount struct {
	ShowTime string
	Quantity string
}

// StockConcession sets the inventory of specified concession to the specified value
func StockConcession(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	printArgs("StockConcession", args)
	if len(args) != 3 {
		msg := fmt.Sprintf("StockConcession(): Incorrect number of arguments. Expecting 3 received %d", len(args))
		return "", errors.New(msg)
	}

	// Use composite key for concessions including theater and concession type
	//E.g.  Regal1-popcorn
	theater := args[0]
	concession := args[1]
	id := fmt.Sprintf("%s-%s", theater, concession)
	con := Concession{theater, concession, args[2], args[3]}

	conJSON, err := json.Marshal(con)
	if err != nil {
		fmt.Println("con to conJSON error: ", err)
		return "", err
	}
	fmt.Println("con to conJSON created: ", string(conJSON))

	err = stub.PutState(id, conJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to stock: %s", args[0])
	}
	results := fmt.Sprintf("Stock created ConID=%s Quantity=%s", con.ConID, con.Inventory)
	return results, nil
}

// ExchangeWaterSoda exchanges 1 water for 1 soda, for the provided Theater and showTime
func ExchangeWaterSoda(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	printArgs("BuyConcession", args)
	if len(args) != 2 {
		msg := fmt.Sprintf("ExchangeWaterSoda(): Incorrect number of arguments. Expecting 2 received %d", len(args))
		return "", errors.New(msg)
	}
	theater := args[0]
	showTime := args[1]
	// first buy buy -1 water then buy 1 soda
	newargs := [4]string{theater, "Water", "-1", ""}
	var myArgs = newargs[:]
	res, err := BuyConcession(stub, myArgs)
	newargs = [4]string{theater, "Soda", "1", showTime}
	myArgs = newargs[:]
	res, err = BuyConcession(stub, myArgs)
	return res, err
}

// BuyConcession purchases 1 or more units (quantity) for a specific concession type
// This reduces the inventory and the purchase transaction is added to the ledger.
func BuyConcession(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	printArgs("BuyConcession", args)

	if len(args) != 4 {
		msg := fmt.Sprintf("BuyConcession(): Incorrect number of arguments. Expecting 4 received %d", len(args))
		return "", errors.New(msg)
	}

	var buycon ConPurchase

	buycon.TxnID = stub.GetTxID()

	// args:  showId, quantity, seller
	buycon.Theater = args[0]
	buycon.ConID = args[1]
	buycon.Quantity = args[2]
	//first, verify that enough tickets are available
	concession, err := ConcessionsAvailable(stub, args)
	if err != nil {
		//not enough available
		fmt.Println(err)
		return "", err
	}

	//Make sure there are enough sodas for this showTime
	if buycon.ConID == "soda" {
		_, err = SodasAvailable(stub, args[2:])
		if err != nil {
			//not enough sodas available for this showtime
			fmt.Println(err)
			return "", err
		}
	}
	// Set current date / time as purchase date
	t := time.Now()
	nowFormatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	buycon.PurchaseTime = nowFormatted //"2019-02-05T13:00:00-05:00"
	//Compute revenue for this sale
	var q, _ = strconv.Atoi(buycon.Quantity)
	var p, _ = strconv.ParseFloat(concession.Price, 64)
	var i, _ = strconv.Atoi(concession.Inventory)
	rev := float64(q) * p
	buycon.Revenue = fmt.Sprintf("%f", rev)
	fmt.Println("BuyCon() - purchase attributes : ", buycon)
	buyconJSON, err := json.Marshal(buycon)
	if err != nil {
		fmt.Println("buycon to buyconJSON error: ", err)
		return "", err
	}
	fmt.Println("buycon to buyconJSON created: ", string(buyconJSON))

	//update the conession with the new inventory count
	concession.Inventory = fmt.Sprintf("%d", i-q)
	conJSON, err := json.Marshal(concession)
	if err != nil {
		fmt.Println("concession to conJSON error: ", err)
		return "", err
	}
	err = stub.PutState(concession.ConID, conJSON)
	if err != nil {
		return "", fmt.Errorf("BuyCon failed to update inventory for: %s", concession.ConID)
	}
	err = stub.PutState(buycon.TxnID, buyconJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to complete ticket purchase for: %s", string(conJSON))
	}
	results := fmt.Sprintf("Buy Concession ConID=%s Quantity=%s Inventory=%s", concession.ConID, buycon.Quantity, concession.Inventory)
	return results, nil
}

// ConcessionsAvailable verifies a quantity of concession units are available for a given type
func ConcessionsAvailable(stub shim.ChaincodeStubInterface, args []string) (Concession, error) {

	var con Concession
	if len(args) != 4 {
		msg := fmt.Sprintf("ConcessionsAvailable(): Incorrect number of arguments. Expecting 3, received %d", len(args))
		return con, errors.New(msg)
	}
	// args:  theater, concession type, quantity
	theater := args[0]
	conID := args[1]
	quantity := args[2]
	showTime := args[3]

	bytes, err := stub.GetState(conID)
	if err != nil {
		return con, fmt.Errorf("ConcessionsAvailable failed to get show information for: %s", conID)
	}

	err = json.Unmarshal(bytes, &con)
	var q, _ = strconv.Atoi(quantity)
	var i, _ = strconv.Atoi(con.Inventory)

	fmt.Printf("ConcessionsAvailable requesting %s of %s. Inventory available is %s", quantity, con.ConID, con.Inventory)

	if q > i {
		return con, fmt.Errorf("Error - Out of Concessions: Attempting to buy %s units of %s at theater %s, but only %s are available", quantity, conID, theater, con.Inventory)
	}

	// Check here for >= 200 sodas for this showTime
	if showTime == "xyz" {
		return con, nil
	}

	return con, nil
}

// SodasAvailable verifies a quantity of sodas are available for a given showTime
func SodasAvailable(stub shim.ChaincodeStubInterface, args []string) (SodaCount, error) {

	var scount SodaCount
	if len(args) != 2 {
		msg := fmt.Sprintf("SodasAvailable(): Incorrect number of arguments. Expecting 2, received %d", len(args))
		return scount, errors.New(msg)
	}
	// args:  theater, concession type, quantity
	quantity := args[0]
	showTime := args[1]

	bytes, err := stub.GetState(showTime)
	if err != nil {
		return scount, fmt.Errorf("SodasAvailable() failed to get current quantity for showTime: %s", showTime)
	}

	err = json.Unmarshal(bytes, &scount)
	curr, _ := strconv.Atoi(scount.Quantity)
	new, _ := strconv.Atoi(quantity)
	total := curr + new
	if total > 200 {
		return scount, fmt.Errorf("Error - Out of Concessions: 200 sodas is the limit for a single showTime. Attempting to buy %s units of sodas for showTime %s, but  %s have been sold", quantity, showTime, scount.Quantity)
	}
	fmt.Printf("Sodas sold for showTime %s is %s, now requesting %s", showTime, scount.Quantity, quantity)

	//update the showTime with the new number of purchased sodas
	quantity = fmt.Sprintf("%d", total)
	scount = SodaCount{showTime, quantity}
	sodaJSON, err := json.Marshal(scount)
	if err != nil {
		fmt.Println("scount to sodaJSON error: ", err)
		return scount, err
	}
	err = stub.PutState(showTime, sodaJSON)
	if err != nil {
		return scount, fmt.Errorf("SodasAvailable failed to update purchased quantity for showTime: %s", showTime)
	}

	return scount, nil
}

var queryString = "{\"selector\": {\"Hall\": \"hall1\"}}"

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
