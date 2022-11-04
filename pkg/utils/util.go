package utils

import (
	"encoding/json"
)

// Recast https://stackoverflow.com/a/72436927
func Recast(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}
