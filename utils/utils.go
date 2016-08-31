package utils

import (
	"encoding/json"
)

func Parse(data []byte) (v interface{}, err error) {
	err = json.Unmarshal(data, &v)
	return
}

func Get(field string, v interface{}) interface{} {
	return v.(map[string]interface{})[field]
}
