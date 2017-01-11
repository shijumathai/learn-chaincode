/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
   "time"
	 "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
const TRAVEL_CONTRACT   = "Paris"
const FEEDBACK_CONTRACT = "Feedback"
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
type contract struct{
	Id			string   `json:"ID"`
		BusinessId  string   `json:"BusinessId"`
		BusinessName string   `json:"BusinessName"`
		Title		string   `json:"Title"`
		Description string   `json:"Description"`
		Conditions  []string `json:"Conditions"`
		Icon        string 	 `json:"Icon"`
		StartDate   time.Time   `json:"StartDate"`
		EndDate		time.Time   `json:"EndDate"`
		Method	    string   `json:"Method"`
		DiscountRate float64  `json:"DiscountRate"`

}
// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err:= stub.PutState("hello_world",[]byte(args[0]))
	if err !=nil {
		return nil,err
	}
	var refnumber int
	refnumber=289907

	 jsonAsBytes, _ := json.Marshal(refnumber)
	err = stub.PutState("refnumber",jsonAsBytes)
	if err !=nil {
		fmt.Println("error creating")
		return nil,err
	}

	var double contract
  	double.Id = TRAVEL_CONTRACT
		double.BusinessId  = "T5940872"
		double.BusinessName = "Open Travel"
		double.Title = "Paris for Less"
		double.Description = "All Paris travel activities are half the stated point price"
		double.Conditions = append(double.Conditions, "Half off dining and travel activities in Paris")
		double.Conditions = append(double.Conditions, "Valid from May 11, 2016")
		double.Icon = ""
		double.Method = "travelContract"

		startDate, _  := time.Parse(time.RFC822, "11 May 16 12:00 UTC")
		double.StartDate = startDate
		endDate, _  := time.Parse(time.RFC822, "31 Dec 60 11:59 UTC")
		double.EndDate = endDate

		jsonAsBytes, _ = json.Marshal(double)
			err = stub.PutState(TRAVEL_CONTRACT, jsonAsBytes)
			if err != nil {
				fmt.Println("Error creating double contract")
				return nil, err
			}


	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}	else if function == "write" {
		return t.write(stub,args)
	}	else if function == "addSmartContract" {											//create a transaction
		return t.addSmartContract(stub, args)
	}
		fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}




func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface,args []string)([]byte,error){
var name, value string
var err error
fmt.Println("running write()")

if len(args) != 2{
return nil ,errors.New ("Incorrect arugments")
}
name=args[0]
value=args[1]
err = stub.PutState(name,[]byte(value))
if err!=nil {
	return nil,err
}
return nil,nil
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err := stub.GetState(name)
		var contractIds []string


		var contract contract
		json.Unmarshal(valAsbytes, &contract)
		//asBytes, _ := json.Marshal(allContracts)

  //  if err != nil {
    //    jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
    //    return nil, errors.New(jsonResp)
  //  }
//json.Unmarshal(valAsbytes,&contract)
	//	asBytes, _ := json.Marshal(allContracts)
//	return valAsbytes, nil
return valAsbytes,nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "read" {                            //read a variable
	        return t.read(stub, args)
	    }
	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
func (t *SimpleChaincode) addSmartContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	// Create new smart contract based on user input
	var smartContract contract

	discountRate, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		smartContract.Title= "Invalid Contract"
	}else{
		smartContract.DiscountRate = discountRate
	}


	smartContract.Id = args[1]
	smartContract.BusinessId  = "T5940872"
	smartContract.BusinessName = "Open Travel"
	smartContract.Title = args[2]
	smartContract.Description = ""
	smartContract.Conditions = append(smartContract.Conditions, args[3])
	smartContract.Conditions = append(smartContract.Conditions, args[4])
	smartContract.Icon = ""
	smartContract.Method = "travelContract"


	jsonAsBytes, _ := json.Marshal(smartContract)
	err = stub.PutState(smartContract.Id, jsonAsBytes)
	if err != nil {
		fmt.Println("Error adding new smart contract")
		return nil, err
	}

	contractIdsAsBytes, _ := stub.GetState("contractIds")
	var contractIds []string
	json.Unmarshal(contractIdsAsBytes, &contractIds)


	var contractIdFound bool
	contractIdFound = false;
	for i := range contractIds{
		if (contractIds[i] == smartContract.Id)  {
			contractIdFound = true;
		}
	}

	if (!contractIdFound) {
		contractIds = append(contractIds, smartContract.Id);
	}


	jsonAsBytes, _ = json.Marshal(contractIds)
	err = stub.PutState("contractIds", jsonAsBytes)
	if err != nil {
		fmt.Println("Error storing contract Ids on blockchain")
		return nil, err
	}

	return nil, nil

}
