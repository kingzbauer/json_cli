package utils

import (
	"reflect"
	"testing"
)

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
