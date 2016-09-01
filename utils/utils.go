package utils

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

var (
	FIELD_SEP string = "."
	INDEX_REG        = regexp.MustCompile(`^\[\d+\]$`)
)

func Parse(data []byte) (v interface{}, err error) {
	err = json.Unmarshal(data, &v)
	return
}

func isIndex(v string) bool {
	return INDEX_REG.Match([]byte(v))
}

func retrieveValueFromIndex(v string) int {
	// remote the opening bracket
	v = strings.Replace(v, "[", "", 1)
	v = strings.Replace(v, "]", "", 1)
	intV, _ := strconv.Atoi(v)
	return intV
}

func get(field string, v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		return t[field]
	case []interface{}:
		if isIndex(field) {
			index := retrieveValueFromIndex(field)
			if index < len(t) {
				return t[index]
			}
		}
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
