// SPDX-License-Identifier: Apache-2.0

package main

/* Imports
* 4 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_employee() - remove a employee from state and from employee index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//// ============================================================================================================================
func delete_employee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_employee")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	authed_by_domain := args[1]

	// get the employee
	employee, err := get_employee(stub, id)
	if err != nil {
		fmt.Println("Failed to find employee by id " + id)
		return shim.Error(err.Error())
	}

	// check authorizing company (see note in set_owner() about how this is quirky)
	if employee.Domain.DomainName != authed_by_domain {
		return shim.Error("The domain '" + authed_by_domain + "' cannot authorize deletion for '" + employee.Domain.DomainName + "'.")
	}

	// remove the employee
	err = stub.DelState(id) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_employee")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Employee - create a new employee, store into chaincode state
//
// ============================================================================================================================
func init_employee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_employee")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	employeeName := args[1]
	designation := args[2]
	dob := args[3]
	domain_id := args[4]
	domain_name := args[5]
	company_id := args[6]
	company_name := args[7]

	// employeeName := strings.ToLower(args[1])
	//designation := args[3]
	//dob := args[4]

	// size, err := strconv.Atoi(args[2])
	// if err != nil {
	// 	return shim.Error("3rd argument must be a numeric string")
	// }

	//check if new doamin exists
	domain, err := get_domain(stub, domain_id)
	if err != nil {
		fmt.Println("Failed to find domain - " + domain_id)
		return shim.Error(err.Error())
	}

	//check authorizing domain (see note in set_owner() about how this is quirky)
	if domain.DomainName != domain_name {
		return shim.Error("The domain '" + domain_name + "' cannot authorize creation for '" + domain.DomainName + "'.")
	}

	//check if new doamin exists
	company, err := get_company(stub, company_id)
	if err != nil {
		fmt.Println("Failed to find domain - " + company_id)
		return shim.Error(err.Error())
	}

	//check authorizing domain (see note in set_owner() about how this is quirky)
	if company.CompanyName != company_name {
		return shim.Error("The company '" + company_name + "' cannot authorize creation for '" + company.CompanyName + "'.")
	}

	//check if employee id already exists
	employee, err := get_employee(stub, id)
	if err == nil {
		fmt.Println("This employee already exists - " + id)
		fmt.Println(employee)
		return shim.Error("This employee already exists - " + id) //all stop a employee by this id exists
	}

	//build the employee json string manually
	// str := `{
	// 	"docType":"employee",
	// 	"id": "` + id + `",
	// 	"color": "` + color + `",
	// 	"size": ` + strconv.Itoa(size) + `,
	// 	"owner": {
	// 		"id": "` + owner_id + `",
	// 		"username": "` + owner.Username + `",
	// 		"company": "` + owner.Company + `"
	// 	}
	// }`

	str := `{
		"docType":"employee", 
		"id": "` + id + `", 
		"employeeName": "` + employeeName + `", 
		"designation": "` + designation + `", 
		"Dob": "` + dob + `", 
		"domain": {
			"id": "` + domain_id + `", 
			
			"domainName": "` + domain.DomainName + `"
		},
		"company": {
			"id": "` + company_id + `", 
			
			"companyName": "` + company.CompanyName + `"
		}
	}`

	err = stub.PutState(id, []byte(str)) //store employee with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_employee")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Company - create a new company , store into chaincode state
//
// ============================================================================================================================
func init_company(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_company")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var company Company
	company.ObjectType = "company"
	company.Id = args[0]
	// company.Username = strings.ToLower(args[1])
	company.CompanyName = args[1]
	company.Enabled = true
	fmt.Println(company)

	//check if user already exists
	_, err = get_company(stub, company.Id)
	if err == nil {
		fmt.Println("This company already exists - " + company.Id)
		return shim.Error("This company already exists - " + company.Id)
	}

	//store company
	companyAsBytes, _ := json.Marshal(company)      //convert to array of bytes
	err = stub.PutState(company.Id, companyAsBytes) //store company by its Id
	if err != nil {
		fmt.Println("Could not store company")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_company employee")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Domain - create a new doamin, store into chaincode state
//
// ============================================================================================================================
func init_domain(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_domain")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var domain Domain
	domain.ObjectType = "domain"
	domain.Id = args[0]
	// owner.Username = strings.ToLower(args[1])
	domain.DomainName = args[1]
	fmt.Println("domain id - " + domain.Id)
	fmt.Println("domain name - " + domain.DomainName)
	domain.Enabled = true
	fmt.Println(domain)

	//check if user already exists
	_, err = get_domain(stub, domain.Id)
	if err == nil {
		fmt.Println("This domain already exists - " + domain.Id)
		return shim.Error("This domain already exists - " + domain.Id)
	}

	//store domain
	domainAsBytes, _ := json.Marshal(domain)      //convert to array of bytes
	err = stub.PutState(domain.Id, domainAsBytes) //store domain by its Id
	if err != nil {
		fmt.Println("Could not store domain")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_domain employee")
	return shim.Success(nil)
}

// ============================================================================================================================
// Set Domain on Employee
//
// Shows off GetState() and PutState()
// ============================================================================================================================
func set_domain(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting set_domain")

	// this is quirky
	// todo - get the "company that authed the transfer" from the certificate instead of an argument
	// should be possible since we can now add attributes to the enrollment cert
	// as is.. this is a bit broken (security wise), but it's much much easier to demo! holding off for demos sake

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var employee_id = args[0]
	var new_domain_id = args[1]
	var authed_by_domain = args[2]
	fmt.Println(employee_id + "->" + new_domain_id + " - |" + authed_by_domain)

	// check if user already exists
	owner, err := get_domain(stub, new_domain_id)
	if err != nil {
		return shim.Error("This domain does not exist - " + new_domain_id)
	}

	// get employee's current state
	employeeAsBytes, err := stub.GetState(employee_id)
	if err != nil {
		return shim.Error("Failed to get employee")
	}
	res := Employee{}
	json.Unmarshal(employeeAsBytes, &res) //un stringify it aka JSON.parse()

	// check authorizing company
	if res.Domain.DomainName != authed_by_domain {
		return shim.Error("The domain '" + authed_by_domain + "' cannot authorize transfers for '" + res.Domain.DomainName + "'.")
	}

	// transfer the employee
	res.Domain.Id = new_domain_id //change the owner
	res.Domain.DomainName = owner.DomainName
	// res.Owner.Company = owner.Company
	jsonAsBytes, _ := json.Marshal(res)       //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the employee with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end set domain")
	return shim.Success(nil)
}

// ============================================================================================================================
// Set Company on Employee
//
// ============================================================================================================================
func set_company(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting set_company")

	// this is quirky
	// todo - get the "company that authed the transfer" from the certificate instead of an argument
	// should be possible since we can now add attributes to the enrollment cert
	// as is.. this is a bit broken (security wise), but it's much much easier to demo! holding off for demos sake

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var employee_id = args[0]
	var new_company_id = args[1]
	var authed_by_company = args[2]
	fmt.Println(employee_id + "->" + new_company_id + " - |" + authed_by_company)

	// check if user already exists
	company, err := get_company(stub, new_company_id)
	if err != nil {
		return shim.Error("This company does not exist - " + new_company_id)
	}

	// get employee's current state
	employeeAsBytes, err := stub.GetState(employee_id)
	if err != nil {
		return shim.Error("Failed to get employee")
	}
	res := Employee{}
	json.Unmarshal(employeeAsBytes, &res) //un stringify it aka JSON.parse()

	// check authorizing company
	if res.Company.CompanyName != authed_by_company {
		return shim.Error("The company '" + authed_by_company + "' cannot authorize transfers for '" + res.Company.CompanyName + "'.")
	}

	// transfer the employee
	res.Company.Id = new_company_id //change the owner
	res.Company.CompanyName = company.CompanyName
	// res.Owner.Company = owner.Company
	jsonAsBytes, _ := json.Marshal(res)       //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the employee with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end set company")
	return shim.Success(nil)
}

// ============================================================================================================================
// Disable Employee Domain
//
// ============================================================================================================================
func disable_domain(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting disable_domain")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var domain_id = args[0]
	var authed_by_domain = args[1]

	// get the employee owner data
	domain, err := get_domain(stub, domain_id)
	if err != nil {
		return shim.Error("This domain does not exist - " + domain_id)
	}

	// check authorizing company
	if domain.DomainName != authed_by_domain {
		return shim.Error("The domain '" + authed_by_domain + "' cannot change another  employee domain")
	}

	// disable the owner
	domain.Enabled = false
	jsonAsBytes, _ := json.Marshal(domain)    //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the owner
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end disable_domain")
	return shim.Success(nil)
}

// ============================================================================================================================
// Disable Employee Company
//
// ============================================================================================================================
func disable_company(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting disable_company")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var company_id = args[0]
	var authed_by_company = args[1]

	// get the employee owner data
	company, err := get_company(stub, company_id)
	if err != nil {
		return shim.Error("This company does not exist - " + company_id)
	}

	// check authorizing company
	if company.CompanyName != authed_by_company {
		return shim.Error("The company '" + authed_by_company + "' cannot change another  employee company")
	}

	// disable the owner
	company.Enabled = false
	jsonAsBytes, _ := json.Marshal(company)   //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the owner
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end disable_company")
	return shim.Success(nil)
}
