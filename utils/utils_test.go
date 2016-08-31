package utils

import (
	"testing"
)

func TestParse(t *testing.T) {
	validJson := `{
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

	_, err := Parse([]byte(validJson))
	if err != nil {
		t.Errorf("Expected to parse content without error. Got \"%v\" instead.", err)
	}

	validJson = `{
        "Id": "d836a5e40aa8974d7076e791ba3c14726bf2dd2cd079652477d6827973969130",
        "Created": "2016-08-31T16:49:33.119587574Z",
        "Path": "/bin/bash",`
	_, err = Parse([]byte(validJson))
	if err == nil {
		t.Errorf("Expected an error while parsing the content.")
	}
}
