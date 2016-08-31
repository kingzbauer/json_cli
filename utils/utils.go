package utils

import (
	"encoding/json"
	"strings"
)

var (
	FIELD_SEP string = "."
)

func Parse(data []byte) (v interface{}, err error) {
	err = json.Unmarshal(data, &v)
	return
}

func get(field string, v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		return t[field]
	}

	return nil
}

func Get(field string, v interface{}) interface{} {
	fields := strings.Split(field, FIELD_SEP)
	result := v

	for _, fieldStr := range fields {
		result = get(fieldStr, result)
		if result == nil {
			return nil
		}
	}
	return result
}
