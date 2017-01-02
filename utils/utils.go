package utils

import (
	"encoding/json"
	"fmt"
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

// ListKeys given a root key, returns all the to level keys under the root key if the
// key exists
func ListKeys(rootKey string, v interface{}) []string {
	if len(rootKey) > 0 {
		v = Get(rootKey, v)
	}

	switch t := v.(type) {
	case []interface{}:
		return arrayKeys(t)
	case map[string]interface{}:
		return mapKeys(t)
	default:
		return nil
	}
}

func mapKeys(m map[string]interface{}) []string {
	v := make([]string, len(m))
	var index int
	for key := range m {
		v[index] = key
		index++
	}

	return v

}

func arrayKeys(arr []interface{}) []string {
	v := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		v[i] = fmt.Sprintf("[%d]", i)
	}

	return v
}

// Search searches for the value for `key` given the data structure upto to `depth` deep
func Search(key string, data interface{}, depth int) interface{} {
	return search(key, data, depth, 0)
}

func search(key string, data interface{}, depth, currentdepth int) interface{} {
	if !isIterable(data) {
		return nil
	}

	v := get(key, data)
	if v != nil {
		return v
	}

	if depth == currentdepth {
		return nil
	}

	var keys []string
	switch t := data.(type) {
	case map[string]interface{}:
		keys = mapKeys(t)
	case []interface{}:
		keys = arrayKeys(t)
	}
	for _, k := range keys {
		v = search(key, get(k, data), depth, currentdepth+1)
		if v != nil {
			return v
		}
	}

	return nil
}

func isIterable(v interface{}) bool {
	switch t := v.(type) {
	case []interface{}, map[string]interface{}:
		_ = t
		return true
	default:
		return false
	}
}
