package utils

import (
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"reflect"
	"strings"
)

// -- json/struct/map helper -- //

// todo:value did not assign to struct...
func JsonToStruct(jsonstring string, struType interface{}) (interface{}, error) {
	t := reflect.TypeOf(struType)
	fmt.Println("struType=", t)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	struType = reflect.New(t)
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	dec.UseNumber()
	err := dec.Decode(&struType)
	return struType, err
}

func StructToJson(stru interface{}) ([]byte, error) {
	return json.Marshal(stru)
}

func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	return mapResult, err
}

func MapToJson(maps ...map[string]interface{}) ([]byte, error) {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return json.Marshal(result)
}

// todo:value did not assign to struct...
func MapToStruct(rawMap map[string]interface{}, struType interface{}) error {
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
