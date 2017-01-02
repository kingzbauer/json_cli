package utils

import (
	"reflect"
	"sort"
	"testing"
)

var MapJSON = `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",
        "Args": [{"bool": true}],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false
        }
     }`

var ArrayJSON = `[
    {"Id": 23,
     "State": {
       "Status": "running",
       "Running": true
     }
    },
    {}
  ]`

type sortableStrings []string

func (s sortableStrings) Len() int {
	return len(s)
}

func (s sortableStrings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableStrings) Less(i, j int) bool {
	return s[i] < s[j]
}

func TestParse(t *testing.T) {
	validJSON := `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",
        "Args": [],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false
        }
     }`

	_, err := Parse([]byte(validJSON))
	if err != nil {
		t.Errorf("Expected to parse content without error. Got \"%v\" instead.", err)
	}

	validJSON = `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",`
	_, err = Parse([]byte(validJSON))
	if err == nil {
		t.Errorf("Expected an error while parsing the content.")
	}
}

func TestGet(t *testing.T) {
	validJSON := `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",
        "Args": [],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false
        }
     }`
	v, _ := Parse([]byte(validJSON))
	expectedV := v.(map[string]interface{})["State"]
	returnedV := Get("State", v)
	if !reflect.DeepEqual(expectedV, returnedV) {
		t.Errorf("Expected: %v, Got: %v", expectedV, returnedV)
	}

	if Get("States", v) != nil {
		t.Errorf("Expected nil for a non existing field")
	}
}

func TestGetNestedFields(t *testing.T) {
	validJSON := `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",
        "Args": [],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false,
            "Log": {
                "Level": "Warn"
            }
        }
     }`

	v, _ := Parse([]byte(validJSON))
	key := "State.Log.Level"
	expectedV := v.(map[string]interface{})["State"].(map[string]interface{})["Log"].(map[string]interface{})["Level"]
	returnedV := Get(key, v)
	if !reflect.DeepEqual(expectedV, returnedV) {
		t.Errorf("Expected %v: Got %v", expectedV, returnedV)
	}

	if Get("State.Log.Nonexisting", v) != nil {
		t.Error("Expected a nil result")
	}
}

func TestCanIndexAnArray(t *testing.T) {
	validJSON := `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",
        "Args": [
            {"id": 56},
            {"id": 23}
        ],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false
        }
     }`
	v, _ := Parse([]byte(validJSON))
	key := "Args.[1].id"
	expectedV := v.(map[string]interface{})["Args"].([]interface{})[1].(map[string]interface{})["id"]
	returnedV := Get(key, v)
	if !reflect.DeepEqual(expectedV, returnedV) {
		t.Errorf("Expected %v: Got %v", expectedV, returnedV)
	}

	// Test for index out of bounds
	key = "Args.[2].id"
	returnedV = Get(key, v)
	if returnedV != nil {
		t.Errorf("Expected a nil value for index out of bounds: Got \"%v\" instead", returnedV)
	}

	// Test for a non-array entry
	key = "[1]"
	returnedV = Get(key, v)
	if returnedV != nil {
		t.Errorf("Expected a nil value for index out of bounds: Got \"%v\" instead", returnedV)
	}

	// Test for a field with the syntax of indexing
	validJSON = `{"[3]": 78}`
	key = "[3]"
	v, _ = Parse([]byte(validJSON))
	expectedV = v.(map[string]interface{})[key]
	returnedV = Get(key, v)
	if !reflect.DeepEqual(expectedV, returnedV) {
		t.Errorf("Expected %v: Got %v", expectedV, returnedV)
	}
}

func TestListKeys(t *testing.T) {
	// add expected keys for the jsonMap
	expectedKeys := sortableStrings([]string{"Id", "Created", "Path", "Args", "State"})
	sort.Sort(expectedKeys)
	parsedV, _ := Parse([]byte(MapJSON))
	returnedKeys := sortableStrings(ListKeys("", parsedV))
	sort.Sort(returnedKeys)
	if !reflect.DeepEqual(expectedKeys, returnedKeys) {
		t.Errorf("Expected %v. Returned %v", expectedKeys, returnedKeys)
	}

	// test for the array json
	expectedKeys = sortableStrings([]string{"[0]", "[1]"})
	sort.Sort(expectedKeys)
	parsedV, _ = Parse([]byte(ArrayJSON))
	returnedKeys = sortableStrings(ListKeys("", parsedV))
	sort.Sort(returnedKeys)
	if !reflect.DeepEqual(expectedKeys, returnedKeys) {
		t.Errorf("Expected %v. Returned %v", expectedKeys, returnedKeys)
	}

	// try a nested key
	expectedKeys = sortableStrings([]string{"Status", "Running", "Paused"})
	sort.Sort(expectedKeys)
	rootKey := "State"
	parsedV, _ = Parse([]byte(MapJSON))
	returnedKeys = sortableStrings(ListKeys(rootKey, parsedV))
	sort.Sort(returnedKeys)
	if !reflect.DeepEqual(expectedKeys, returnedKeys) {
		t.Errorf("Expected %v. Returned %v", expectedKeys, returnedKeys)
	}
}

func TestSearchKey(t *testing.T) {
	// Test a key search in an array
	key := "bool"
	expectedV := true
	// a depth of 0 should miss the key
	parsedV, _ := Parse([]byte(MapJSON))
	receivedV := Search(key, parsedV, 0)
	if reflect.DeepEqual(receivedV, expectedV) {
		t.Errorf("Expected nil, got %#v", receivedV)
	}

	receivedV = Search(key, parsedV, 2)
	if !reflect.DeepEqual(receivedV, expectedV) {
		t.Errorf("Expected %#v, got %#v", expectedV, receivedV)
	}

	// Test a key search in a map
	key = "Status"
	statusV := "running"
	receivedV = Search(key, parsedV, 1)
	if !reflect.DeepEqual(statusV, receivedV) {
		t.Errorf("Expected %#v, got %#v", statusV, receivedV)
	}
}
