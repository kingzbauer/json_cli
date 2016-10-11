package utils

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

var (
	// FieldSep defines what is used to separate fields
	FieldSep = "."
	// IndexReg defines the regular expression that matches an array indexing field
	IndexReg = regexp.MustCompile(`^\[\d+\]$`)
)

// Parse unmarshals a byte string to its respective Go data structure
func Parse(data []byte) (v interface{}, err error) {
	err = json.Unmarshal(data, &v)
	return
}

func isIndex(v string) bool {
	return IndexReg.Match([]byte(v))
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

// Get returns the value of a given key from the data present in v.
// You can dig deep into the tree by separating field names with a period `.`
// Indexing an array can be done via the normal indexing syntax e.g `[0]` so that a whole
// key could be something like this: `key1.[0].key2`
func Get(field string, v interface{}) interface{} {
	fields := strings.Split(field, FieldSep)
	result := v

	for _, fieldStr := range fields {
		result = get(fieldStr, result)
		if result == nil {
			return nil
		}
	}
	return result
}

func ListKeys(rootKey string, v interface{}) []string {
	return nil
}
