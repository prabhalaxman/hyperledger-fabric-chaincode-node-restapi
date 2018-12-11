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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key) //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes) //send it onward
}

func read_all_data(stub shim.ChaincodeStubInterface) pb.Response {
	type AllData struct {
		Employee []Employee `json:"employee"`
		Domain   []Domain   `json:"domain"`
		Company  []Company  `json:"company"`
	}
	var allData AllData

	// ---- Get All Employees ---- //
	employeeIterator, err := stub.GetStateByRange("e0", "e9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer employeeIterator.Close()

	for employeeIterator.HasNext() {
		aKeyValue, err := employeeIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on employee id - ", queryKeyAsStr)
		var employee Employee
		json.Unmarshal(queryValAsBytes, &employee)            //un stringify it aka JSON.parse()
		allData.Employee = append(allData.Employee, employee) //add this employee to the list
	}
	fmt.Println("employee array - ", allData.Employee)

	// ---- Get All domain ---- //
	doaminIterator, err := stub.GetStateByRange("d0", "d9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer doaminIterator.Close()

	for doaminIterator.HasNext() {
		aKeyValue, err := doaminIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on domain id - ", queryKeyAsStr)
		var domain Domain
		json.Unmarshal(queryValAsBytes, &domain) //un stringify it aka JSON.parse()

		if domain.Enabled { //only return enabled owners
			allData.Domain = append(allData.Domain, domain) //add this employee to the list
		}
	}
	fmt.Println("domain array - ", allData.Domain)

	// ---- Get All domain ---- //
	companyIterator, err := stub.GetStateByRange("c0", "c9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer companyIterator.Close()

	for companyIterator.HasNext() {
		aKeyValue, err := companyIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on company id - ", queryKeyAsStr)
		var company Company
		json.Unmarshal(queryValAsBytes, &company) //un stringify it aka JSON.parse()

		if company.Enabled { //only return enabled owners
			allData.Company = append(allData.Company, company) //add this employee to the list
		}
	}
	fmt.Println("company array - ", allData.Company)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(allData) //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Get history of asset
//
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId  string   `json:"txId"`
		Value Employee `json:"value"`
	}
	var history []AuditHistory
	var employee Employee

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	employeeId := args[0]
	fmt.Printf("- start getHistoryForEmployee: %s\n", employeeId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(employeeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                   //copy transaction id over
		json.Unmarshal(historyData.Value, &employee) //un stringify it aka JSON.parse()
		if historyData.Value == nil {                //employee has been deleted
			var emptyEmployee Employee
			tx.Value = emptyEmployee //copy nil employee
		} else {
			json.Unmarshal(historyData.Value, &employee) //un stringify it aka JSON.parse()
			tx.Value = employee                          //copy employee over
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForEmployee returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of asset - performs a range query based on the start and end keys provided.
//
// ============================================================================================================================
func getAllEmployee(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryResultKey := aKeyValue.Key
		queryResultValue := aKeyValue.Value

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getEmployeeByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func getEmployeeById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	employeeAsBytes, _ := stub.GetState(args[0])
	if employeeAsBytes == nil {
		return shim.Error("Could not get employee")
	}
	return shim.Success(employeeAsBytes)

}
