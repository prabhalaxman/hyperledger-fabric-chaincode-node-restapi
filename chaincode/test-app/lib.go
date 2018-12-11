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
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Get Employees - get a employees asset from ledger
// ============================================================================================================================
func get_employee(stub shim.ChaincodeStubInterface, id string) (Employee, error) {
	var employee Employee
	employeeAsBytes, err := stub.GetState(id) //getState retreives a key/value from the ledger
	if err != nil {                           //this seems to always succeed, even if key didn't exist
		return employee, errors.New("Failed to find employee - " + id)
	}
	json.Unmarshal(employeeAsBytes, &employee) //un stringify it aka JSON.parse()

	if employee.Id != id { //test if employee is actually here or just nil
		return employee, errors.New("Employee does not exist - " + id)
	}

	return employee, nil
}

// ============================================================================================================================
// Get Domain - get the domain asset from ledger
// ============================================================================================================================
func get_domain(stub shim.ChaincodeStubInterface, id string) (Domain, error) {
	var doamin Domain
	doaminAsBytes, err := stub.GetState(id) //getState retreives a key/value from the ledger
	if err != nil {                         //this seems to always succeed, even if key didn't exist
		return doamin, errors.New("Failed to get doamin - " + id)
	}
	json.Unmarshal(doaminAsBytes, &doamin) //un stringify it aka JSON.parse()

	if len(doamin.DomainName) == 0 { //test if domain is actually here or just nil
		return doamin, errors.New("Doamin does not exist - " + id + ", '" + doamin.DomainName + "'")
	}

	return doamin, nil
}

// ============================================================================================================================
// Get Company - get the company asset from ledger
// ============================================================================================================================
func get_company(stub shim.ChaincodeStubInterface, id string) (Company, error) {
	var company Company
	companyAsBytes, err := stub.GetState(id) //getState retreives a key/value from the ledger
	if err != nil {                          //this seems to always succeed, even if key didn't exist
		return company, errors.New("Failed to get company - " + id)
	}
	json.Unmarshal(companyAsBytes, &company) //un stringify it aka JSON.parse()

	if len(company.CompanyName) == 0 { //test if company is actually here or just nil
		return company, errors.New("company does not exist - " + id + ", '" + company.CompanyName + "'")
	}

	return company, nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error {
	for i, val := range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}
