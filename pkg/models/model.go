package models

import (
	_ "github.com/bdarge/api-gateway/cmd/docs"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type Model struct {
	ID        uint32    `json:"id"` // https://stackoverflow.com/a/21152548
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

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

type DispositionItem struct {
	Model
	Description string  `json:"description"`
	Qty         uint32  `json:"qty"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unitPrice"`
} // @name DispositionItem

type Disposition struct {
	Model
	Currency     string            `json:"currency"`
	Description  string            `json:"description" binding:"required"`
	DeliveryDate time.Time         `json:"deliveryDate" binding:"required"`
	CustomerId   uint32            `json:"customerId" binding:"required"`
	CreatedBy    uint32            `json:"createdBy" binding:"required"`
	RequestType  string            `json:"requestType" binding:"required,oneof=order quote"`
	Items        []DispositionItem `json:"items"`
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

type CreateResponse struct {
	ID uint32 `json:"id"`
} //@name DispositionResponse

// swagger:parameters get_disposition
type _ struct {
	// The ID of a disposition
	// in:path
	ID string `json:"id"`
}

// swagger:parameters get_dispositions
type DispositionsRequest struct {
	// Page
	// in:query
	Page uint32 `json:"page"`
	// Limit (max 100)
	// in:query
	Limit uint32 `json:"limit"`
	// in:query
	RequestType string `json:"requestType"`
	// in:query
	Search string `json:"search"`
}

type Dispositions struct {
	Total int32         `json:"total" format:"int32"`
	Page  int32         `json:"page"  format:"int32"`
	Limit int32         `json:"limit" format:"int32"`
	Data  []Disposition `json:"data"`
} // @name Dispositions

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
} // @name ErrorResponse

type ErrorResponse400 struct {
	Errors []ErrorMsg `json:"errors"`
} // @name ErrorResponse400

type Customer struct {
	Model
	Email        string        `json:"email" binding:"required"`
	Name         string        `json:"name" binding:"required"`
	Dispositions []Disposition `json:"dispositions"`
}
