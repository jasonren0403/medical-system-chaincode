package utils

import (
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"reflect"
)

// UnknownTransactionHandler returns a shim error
// with details of a bad transaction request
func UnknownTransactionHandler(ctx contractapi.TransactionContextInterface) error {
	fcn, args := ctx.GetStub().GetFunctionAndParameters()
	return fmt.Errorf("invalid function %s passed with args %v", fcn, args)
}

func ReadInitDataFromFile(filename string) (interface{}, error) {
	return nil, nil
}

// -- json/struct/map helper -- //

func JsonToStruct(jsonstring string, struType struct{}) (struct{}, error) {
	err := json.Unmarshal([]byte(jsonstring), &struType)
	return struType, err
}

func StructToJson(stru struct{}) ([]byte, error) {
	return json.Marshal(stru)
}

func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	return mapResult, err
}

func MapToJson(baseInstance []map[string]interface{}, mapInstances ...map[string]interface{}) ([]byte, error) {
	for _, v := range mapInstances {
		baseInstance = append(baseInstance, v)
	}
	return json.Marshal(baseInstance)
}

func MapToStruct(rawMap map[string]interface{}, struType struct{}) error {
	return mapstructure.Decode(rawMap, &struType)
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
