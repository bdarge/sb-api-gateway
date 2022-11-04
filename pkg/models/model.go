package models

import (
	_ "github.com/bdarge/sb-api-gateway/cmd/docs"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string
	Error  string
}

type LoginResponse struct {
	Response
	Token string
}

// Account to register
type Account struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
} // @name Account

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} // @name Login

type Disposition struct {
	Currency     string `json:"currency"`
	Description  string `json:"description" binding:"required"`
	DeliveryDate string `json:"deliveryDate" binding:"required"`
	CustomerId   int64  `json:"customerId" binding:"required"`
	CreatedBy    int64  `json:"createdBy" binding:"required"`
	RequestType  string `json:"requestType" binding:"required,oneof=order quote"`
} // @name Disposition

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

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

type DispositionResponse struct {
	ID int64 `json:"id"`
} //@name DispositionResponse

// swagger:parameters get_disposition
type _ struct {
	// The ID of a disposition
	// in:path
	ID string `json:"id"`
}

// swagger:parameters get_dispositions
type _ struct {
	// Page
	// in:query
	Page int
	// Limit (max 100)
	// in:query
	Limit int
	// in:query
	RequestTye string
}

type Dispositions struct {
	Total int           `json:"total" format:"int64"`
	Page  int           `json:"page"  format:"int64"`
	Limit int           `json:"limit" format:"int64"`
	Data  []Disposition `json:"data"`
} // @name Dispositions

type ErrorResponse struct {
	Error   string `json:"error""`
	Message string `json:"message""`
} // @name ErrorResponse

type ErrorResponse400 struct {
	Errors []ErrorMsg `json:"errors""`
} // @name ErrorResponse400
