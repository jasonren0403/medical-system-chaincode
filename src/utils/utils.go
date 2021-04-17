package utils

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UnknownTransactionHandler returns a shim error
// with details of a bad transaction request
func UnknownTransactionHandler(ctx contractapi.TransactionContextInterface) error {
	fcn, args := ctx.GetStub().GetFunctionAndParameters()
	return fmt.Errorf("invalid function %s passed with args %v", fcn, args)
}
