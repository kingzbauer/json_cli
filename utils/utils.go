package utils

import (
	"encoding/json"
)

func Parse(data []byte) (v interface{}, err error) {
	err = json.Unmarshal(data, &v)
	return
}
