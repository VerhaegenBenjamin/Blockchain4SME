/*
 * The sample smart contract for documentation topic:
 * cross border funds transfer
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the bank structure, with 3 properties.  Structure tags are used by encoding/json library
type Bank struct {
	Name     string  `json:"name"`
	BankID   string  `json:"bankID"`
	Country  string  `json:"country"`
	Currency string  `json:"currency"`
	Reserves float64 `json:"reserves"`
}

// Define the customer structure, with 3 properties.  Structure tags are used by encoding/json library
type Customer struct {
	Name           string  `json:"name"`
	CustID         string  `json:"custID"`
	Country        string  `json:"country"`
	Currency       string  `json:"currency"`
	Balance        float64 `json:"balance"`
	CustomerBankID string  `json:"customerBankID"`
}


 * The Init method is called when the Smart Contract "banks" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract ""
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAll" { //return all the assets on the ledger
		return s.queryAll(APIstub, args)
	} else if function == "query" { //single bank or customer or forexPair
		return s.query(APIstub, args)
	} else if function == "createBank" {
		return s.createBank(APIstub, args)
	} else if function == "createCustomer" {
		return s.createCustomer(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/** ----------------------------------------------------------------------**/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	banks := []Bank{
		{Name: "US_Bank", BankID: "US_Bank", Country: "USA", Currency: "USD", Reserves: 1000000.0},
		{Name: "UK_Bank", BankID: "UK_Bank", Country: "UK ", Currency: "GBP", Reserves: 1000000.0},
		{Name: "Japan_Bank", BankID: "Japan_Bank", Country: "JAPAN", Currency: "JPY", Reserves: 10000000.0},
	}

	customers := []Customer{
		{Name: "US_John_Doe", CustID: "123", Country: "US", Currency: "USD", Balance: 10000.0, CustomerBankID: "US_Bank"},
		{Name: "US_Alice", CustID: "456", Country: "US", Currency: "USD", Balance: 10000.0, CustomerBankID: "US_Bank"},
		{Name: "UK_John_Doe", CustID: "123", Country: "UK", Currency: "GBP", Balance: 10000.0, CustomerBankID: "UK_Bank"},
		{Name: "UK_Alice", CustID: "456", Country: "UK", Currency: "GBP", Balance: 10000.0, CustomerBankID: "UK_Bank"},
		{Name: "JPY_John_Doe", CustID: "123", Country: "Japan", Currency: "JPY", Balance: 1000000.0, CustomerBankID: "Japan_Bank"},
		{Name: "JPY_Alice", CustID: "456", Country: "Japan", Currency: "JPY", Balance: 1000000.0, CustomerBankID: "Japan_Bank"},
	}

	writeCustomerToLedger(APIstub, customers)
	writeBankToLedger(APIstub, banks)

	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func writeBankToLedger(APIStub shim.ChaincodeStubInterface, banks []Bank) sc.Response {
//Fill in WriteBankToLedger Functionality

	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func writeCustomerToLedger(APIStub shim.ChaincodeStubInterface, customers []Customer) sc.Response {
	//Fill in Write CustomerToLedger Functionality
	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func (s *SmartContract) queryAll(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments for querying all assets. Expecting 1")
	}
	//collection := args[0]
	//startKey := collection + "0"
	//endKey := collection + "99"

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString("\n,")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("\n]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/** ----------------------------------------------------------------------**/

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	asBytes, _ := APIstub.GetState(args[0])
	return shim.Success(asBytes)
}



/** ----------------------------------------------------------------------------------------------
ceate bank needs 5 args
Name     string  `json:"name"`
BankID	 string  `json:"bankID"`
Country  string  `json:"country"`
Currency string  `json:"currency"`
Reserves float64 `json:"reserves"`
args: ['EU_Bank', 'EU_Bank', 'Europe', 'EURO', '1000000.0'],
*/
func (s *SmartContract) createBank(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//*fill in the createBank Functionality
	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
createCustomer needs 6 args
	Name     string  `json:"name"`
	CustID   string  `json:"custID"`
	Country  string  `json:"country"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	CustomerBankID string `json:"customerBankID"`
	["US_Mary_Jane", "789",  "US", "USD", 100000.0, "US_Bank"],
*/
func (s *SmartContract) createCustomer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
  // Fill in the CreateCustomer Functionality

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

	fmt.Println("successfully initialized smart contract")

}
