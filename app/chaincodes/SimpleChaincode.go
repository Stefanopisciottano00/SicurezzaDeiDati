package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SimpleChaincode is a simple chaincode implementation
type SimpleChaincode struct {
	contractapi.Contract
}

// Counter structure to hold the state
type Counter struct {
	Value int `json:"value"`
}

// Init initializes the chaincode
func (s *SimpleChaincode) Init(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// Increment increments the counter value by 1
func (s *SimpleChaincode) Increment(ctx contractapi.TransactionContextInterface) error {
	counter, err := s.getCounter(ctx)
	if err != nil {
		return err
	}

	counter.Value++
	return s.saveCounter(ctx, counter)
}

// GetValue retrieves the current value of the counter
func (s *SimpleChaincode) GetValue(ctx contractapi.TransactionContextInterface) (*Counter, error) {
	return s.getCounter(ctx)
}

// getCounter retrieves the counter from the world state
func (s *SimpleChaincode) getCounter(ctx contractapi.TransactionContextInterface) (*Counter, error) {
	existing, err := ctx.GetStub().GetState("counter")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if existing == nil {
		return &Counter{Value: 0}, nil
	}

	var counter Counter
	err = json.Unmarshal(existing, &counter)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize JSON: %v", err)
	}

	return &counter, nil
}

// saveCounter saves the counter to the world state
func (s *SimpleChaincode) saveCounter(ctx contractapi.TransactionContextInterface, counter *Counter) error {
	counterJSON, err := json.Marshal(counter)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("counter", counterJSON)
}

func main() {
	simpleChaincode, err := contractapi.NewChaincode(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
		return
	}

	if err := simpleChaincode.Start(); err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
	}
}
