package utils

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Recast https://stackoverflow.com/a/72436927
func Recast(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}

// GetErrorMsg parse error message
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "oneof":
		var params = strings.Split(fe.Param(), " ")
		var result = "Should be one of the following: "
		for index, p := range params {
			if index == len(params)-1 {
				result += ", or '" + p + "'"
			} else if index == 0 {
				result += "'" + p + "'"
			} else {
				result += ", '" + p + "'"
			}
		}
		return result
	}
	return "Unknown error"
}
