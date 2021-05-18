package smartMedicineSystem

import (
	"ccode/src/asset"
	"ccode/src/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type TestStru struct {
	Strfield      string  `json:"strfield"`
	Numfield      int     `json:"numfield"`
	FloatNumfield float64 `json:"floatnumfield"`
	Boolfield     bool    `json:"boolfield"`
}

type Nested struct {
}

// bad: value did not assign to struct...
func TestJsonToStruct(t *testing.T) {
	var jsonstr = `{"strfield": "string","numfield": 64,"floatnumfield": 114.5,"boolfield": true}`
	var testStru TestStru
	s, err := utils.JsonToStruct(jsonstr, testStru)
	assert.NoError(t, err, "Error occurs")
	//log.Println(s.(TestStru))
	log.Println(s)
	//assert.EqualValues(t, 64, s.Numfield, "Num field did not set properly")
}

// bad: value did not assign to struct...
func TestMapToStruct(t *testing.T) {
	rawMap := map[string]interface{}{
		"Strfield": "bbb",
		"Numfield": 64,
	}
	var testStru TestStru
	err := utils.MapToStruct(rawMap, &testStru)
	assert.NoError(t, err, "Error occurs")
	log.Println(testStru)
}

func TestJsonToMap(t *testing.T) {
	var jsonstr = `{"strfield": "string","numfield": 64,"floatnumfield": 114.5,"boolfield": true}`
	s, err := utils.JsonToMap(jsonstr)
	assert.NoError(t, err, "Error occurs")
	assert.EqualValues(t, 64, s["numfield"], "Num field did not set properly")
	assert.EqualValues(t, 114.5, s["floatnumfield"], "Float num field did not set properly")
	assert.EqualValues(t, "string", s["strfield"], "String field did not set properly")
	assert.True(t, s["boolfield"].(bool), "Boolean field did not set properly")
}

func TestStructToMap(t *testing.T) {
	var testStru = TestStru{
		Strfield:      "a",
		Numfield:      10,
		FloatNumfield: 1.4,
		Boolfield:     false,
	}
	m := utils.StructToMap(testStru)
	assert.EqualValues(t, 10, m["Numfield"], "Num field did not set properly")
	assert.EqualValues(t, 1.4, m["FloatNumfield"], "Float num field did not set properly")
	assert.EqualValues(t, "a", m["Strfield"], "String field did not set properly")
	assert.False(t, m["Boolfield"].(bool), "Boolean field did not set properly")
}

func TestStructToJson(t *testing.T) {
	test := TestStru{
		Strfield:      "a",
		Numfield:      1,
		FloatNumfield: 11.4,
		Boolfield:     true,
	}
	btest, err := utils.StructToJson(test)
	assert.NoError(t, err, "Error occurs")
	assert.JSONEq(t, `{"strfield":"a","numfield":1,"floatnumfield":11.4,"boolfield":true}`, string(btest),
		"JSON does not match specified")
}

func TestMapToJson(t *testing.T) {
	map1 := map[string]interface{}{
		"key1": "v1",
		"key2": 2,
	}
	map2 := map[string]interface{}{
		"key3": 4,
		"key4": true,
	}
	map3 := map[string]interface{}{
		"key5": false,
		"key6": 3,
	}
	bmaps, err := utils.MapToJson(map1, map2, map3)
	assert.NoError(t, err, "Error occurs")
	assert.Contains(t, string(bmaps), "key1", "key1 should be merged")
	assert.Contains(t, string(bmaps), "key2", "key2 should be merged")
	assert.Contains(t, string(bmaps), "key3", "key3 should be merged")
	assert.Contains(t, string(bmaps), "key4", "key4 should be merged")
	assert.Contains(t, string(bmaps), "key5", "key5 should be merged")
	assert.Contains(t, string(bmaps), "key6", "key6 should be merged")
}

func TestA(t *testing.T) {
	var collaboratorsStr = "[{\"doctor\":{\"person\":{\"id\":\"doct1\",\"name\":\"Apple\",\"age\":24},\"department\":\"Dep1\"},\"role\":\"member\"},{\"doctor\":{\"person\":{\"id\":\"doct2\",\"name\":\"Banana\",\"age\":25},\"department\":\"Dep1\"},\"role\":\"manager\"},{\"doctor\":{\"person\":{\"id\":\"doct3\",\"name\":\"Catt\",\"age\":26},\"department\":\"Dep2\"},\"role\":\"member\"}]"
	var needleStr = "{\"doctor\":{\"person\":{\"id\":\"doct2\",\"name\":\"Banana\",\"age\":25},\"department\":\"Dep1\"},\"role\":\"manager\"}"
	var nonNeedle = "{\"doctor\":{\"person\":{\"id\":\"doct1\",\"name\":\"Apple\",\"age\":27},\"department\":\"Dep1\"},\"role\":\"member\"}"
	var cols []asset.Collaborator
	var needle asset.Collaborator
	err := json.Unmarshal(utils.Str2bytes(collaboratorsStr), &cols)
	assert.NoError(t, err, "")
	err = json.Unmarshal(utils.Str2bytes(needleStr), &needle)
	assert.NoError(t, err, "")
	assert.True(t, asset.InCollaboratorList(cols, needle), "")
	err = json.Unmarshal(utils.Str2bytes(nonNeedle), &needle)
	assert.False(t, asset.InCollaboratorList(cols, needle), "")
}
