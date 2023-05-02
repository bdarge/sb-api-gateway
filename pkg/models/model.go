package models

import (
	_ "github.com/bdarge/api-gateway/cmd/docs"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type Model struct {
	ID        uint32     `json:"id"` // https://stackoverflow.com/a/21152548
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
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

type TransactionItem struct {
	Model
	Description string  `json:"description"`
	Qty         uint32  `json:"qty"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unitPrice"`
} // @name TransactionItem

type NewTransaction struct {
	Model
	Currency     string            `json:"currency"`
	Description  string            `json:"description" binding:"required"`
	DeliveryDate time.Time         `json:"deliveryDate" binding:"required"`
	CustomerId   uint32            `json:"customerId" binding:"required"`
	CreatedBy    uint32            `json:"createdBy" binding:"required"`
	RequestType  string            `json:"requestType" binding:"required,oneof=order quote"`
	Items        []TransactionItem `json:"items"`
} // @name Transaction

// swagger:parameters update_transaction
type UpdateTransaction struct {
	ID           uint32     `json:"id"`
	Currency     string     `json:"currency"`
	Description  string     `json:"description"`
	DeliveryDate *time.Time `json:"deliveryDate"`
	CustomerId   uint32     `json:"customerId"`
	CreatedBy    uint32     `json:"createdBy"`
	RequestType  string     `json:"requestType" binding:"oneof=order quote ''"`
}

type Transaction struct {
	Model
	Currency     string            `json:"currency"`
	Description  string            `json:"description" binding:"required"`
	DeliveryDate time.Time         `json:"deliveryDate" binding:"required"`
	CustomerId   uint32            `json:"customerId" binding:"required"`
	CreatedBy    uint32            `json:"createdBy" binding:"required"`
	RequestType  string            `json:"requestType" binding:"required,oneof=order quote"`
	Customer     Customer          `json:"customer" binding:"not_required"`
	Items        []TransactionItem `json:"items"`
} // @name Transaction

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
} //@name TransactionResponse

// swagger:parameters get_transaction
type _ struct {
	// The ID of a transaction
	// in:path
	ID string `json:"id"`
}

// swagger:parameters get_transactions
type TransactionsRequest struct {
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
	// in:query
	SortProperty string `json:"sortProperty"`
	// in:query
	SortDirection string `json:"sortDirection"`
}

type Transactions struct {
	Total int32         `json:"total" format:"int32"`
	Page  int32         `json:"page"  format:"int32"`
	Limit int32         `json:"limit" format:"int32"`
	Data  []Transaction `json:"data"`
} // @name Transactions

// swagger:parameters get_customers
type CustomersRequest struct {
	// Page
	// in:query
	Page uint32 `json:"page"`
	// Limit (max 100)
	// in:query
	Limit uint32 `json:"limit"`
	// in:query
	Search string `json:"search"`
	// in:query
	SortProperty string `json:"sortProperty"`
	// in:query
	SortDirection string `json:"sortDirection"`
}

type Customers struct {
	Total int32      `json:"total" format:"int32"`
	Page  int32      `json:"page"  format:"int32"`
	Limit int32      `json:"limit" format:"int32"`
	Data  []Customer `json:"data"`
} // @name Customers

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
} // @name ErrorResponse

type ErrorResponse400 struct {
	Errors []ErrorMsg `json:"errors"`
} // @name ErrorResponse400

type Customer struct {
	Model
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

// swagger:parameters update_customer
type UpdateCustomer struct {
	ID    uint32 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Role struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

// Address Model
type Address struct {
	Model
	Street        string `json:"street"`
	PostalCode    string `json:"postalCode"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	LandLinePhone string `json:"landlinePhone"`
	MobilePhone   string `json:"mobilePhone"`
	UserID        uint32 `json:"userId"`
	BusinessID    uint32 `json:"businessId"`
}

// Business Model.
type Business struct {
	Model
	Name       string  `json:"name"`
	HourlyRate uint32  `gorm:"column:hourly_rate" json:"hourlyRate"`
	Vat        uint32  `json:"vat"`
	Address    Address `json:"address"`
}

// User Model
type User struct {
	Model
	Username     string        `json:"username"`
	Address      Address       `json:"address"`
	Transactions []Transaction `json:"transactions"`
	Roles        []Role        `json:"roles"`
	BusinessID   uint32        `json:"businessId"`
	Business     Business      `json:"business"`
}

// swagger:parameters update_user
type UpdateUser struct {
	ID            uint32 `json:"id"`
	UserName      string `json:"username"`
	BusinessName  string `json:"businessName"`
	Street        string `json:"street"`
	PostalCode    string `json:"postalCode"`
	City          string `json:"city"`
	Country       string `json:"country"`
	LandLinePhone string `json:"landlinePhone"`
	MobilePhone   string `json:"mobilePhone"`
}
