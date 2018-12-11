// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

/* Imports
* 4 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Employee struct {
	ObjectType   string          `json:"docType"` //field for couchdb
	Id           string          `json:"id"`
	EmployeeName string          `json:"employeeName"`
	Designation  string          `json:"designation"`
	Dob          string          `json:"dob"`
	Domain       DomainRelation  `json:"domain"`
	Company      CompanyRelation `json:"company"`
}

type Domain struct {
	ObjectType string `json:"docType"` //field for couchdb
	Id         string `json:"id"`
	DomainName string `json:"domainName"`
	Enabled    bool   `json:"enabled"` //disabled owners will not be visible to the application
}

type Company struct {
	ObjectType  string `json:"docType"` //field for couchdb
	Id          string `json:"id"`
	CompanyName string `json:"companyName"`
	Enabled     bool   `json:"enabled"` //disabled owners will not be visible to the application
}

type DomainRelation struct {
	Id         string `json:"id"`
	DomainName string `json:"domainName"` //this is mostly cosmetic/handy, the real relation is by Id not Username

}

type CompanyRelation struct {
	Id          string `json:"id"`
	CompanyName string `json:"companyName"`
}

/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Test Fabric Is Starting Up")
	funcName, args := stub.GetFunctionAndParameters()
	var number int
	var err error
	txId := stub.GetTxID()

	fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)

	// expecting 1 arg for instantiate or upgrade
	if len(args) == 1 {
		fmt.Println("  GetFunctionAndParameters() arg[0] length", len(args[0]))

		// expecting arg[0] to be length 0 for upgrade
		if len(args[0]) == 0 {
			fmt.Println("  Uh oh, args[0] is empty...")
		} else {
			fmt.Println("  Great news everyone, args[0] is not empty")

			// convert numeric string to integer
			number, err = strconv.Atoi(args[0])
			if err != nil {
				return shim.Error("Expecting a numeric string argument to Init() for instantiate")
			}

			// this is a very simple test. let's write to the ledger and error out on any errors
			// it's handy to read this right away to verify network is healthy if it wrote the correct value
			err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
			if err != nil {
				return shim.Error(err.Error()) //self-test fail
			}
		}
	}

	// showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

	// store compatible employee application version
	err = stub.PutState("employee_ui", []byte("4.0.1"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ready for action") //self-test pass
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "read" { //generic read ledger
		return read(stub, args)
	} else if function == "write" { //generic writes to ledger
		return write(stub, args)
	} else if function == "delete_employee" { //deletes a employee from state
		return delete_employee(stub, args)
	} else if function == "init_employee" { //create a new employee
		return init_employee(stub, args)
	} else if function == "set_domain" { //change domain of a employee
		return set_domain(stub, args)
	} else if function == "set_company" { //change company of a employee
		return set_company(stub, args)
	} else if function == "init_domain" { //create a new domain
		return init_domain(stub, args)
	} else if function == "init_company" { //create a new company
		return init_company(stub, args)
	} else if function == "read_all_data" { //read everything, (employee + domain + company)
		return read_all_data(stub)
	} else if function == "getHistory" { //read history of a employee (audit)
		return getHistory(stub, args)
	} else if function == "getAllEmployee" { //read a all employees by start and stop id
		return getAllEmployee(stub, args)
	} else if function == "getEmployeeById" { //read a employee by id
		return getEmployeeById(stub, args)
	} else if function == "disable_domain" { //disable a domain
		return disable_domain(stub, args)
	} else if function == "disable_company" { //disable a company
		return disable_company(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
