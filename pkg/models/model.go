package models

import (
	"encoding/json"
	"time"
)

// Model basic Model
type Model struct {
	ID        uint32     `json:"id"` // https://stackoverflow.com/a/21152548
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// Response Model
type Response struct {
	Status string
	Error  string
}

// LoginResponse Model
type LoginResponse struct {
	Response
	Token string
}


// Login Model
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} // @name Login

// TransactionItem Model
type TransactionItem struct {
	Model
	Description   string      `json:"description"`
	Qty           json.Number `json:"qty"`
	Unit          string      `json:"unit"`
	UnitPrice     json.Number `json:"unitPrice"`
	TransactionID uint32      `json:"transactionId"`
} // @name TransactionItem

// NewTransaction Model
type NewTransaction struct {
	Model
	Currency     string            `json:"currency"`
	Description  string            `json:"description" binding:"required"`
	DeliveryDate time.Time         `json:"deliveryDate" binding:"required"`
	CustomerID   uint32            `json:"customerId" binding:"required"`
	CreatedBy    uint32            `json:"createdBy" binding:"required"`
	RequestType  string            `json:"requestType" binding:"required,oneof=order quote"`
	Items        []TransactionItem `json:"items"`
} // @name Transaction

// UpdateTransaction Model swagger:parameters update_transaction
type UpdateTransaction struct {
	ID           uint32     `json:"id"`
	Currency     string     `json:"currency"`
	Description  string     `json:"description"`
	DeliveryDate *time.Time `json:"deliveryDate"`
	CustomerID   uint32     `json:"customerId"`
	CreatedBy    uint32     `json:"createdBy"`
	RequestType  string     `json:"requestType" binding:"oneof=order quote ''"`
}

// Transaction Model
type Transaction struct {
	Model
	Currency     string            `json:"currency"`
	Description  string            `json:"description" binding:"required"`
	DeliveryDate time.Time         `json:"deliveryDate" binding:"required"`
	CreatedBy    uint32            `json:"createdBy" binding:"required"`
	RequestType  string            `json:"requestType" binding:"required,oneof=order quote"`
	Customer     Customer          `json:"customer" binding:"not_required"`
	Items        []TransactionItem `json:"items"`
} // @name Transaction

// ErrorMsg Model
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// CreateResponse Model
type CreateResponse struct {
	ID uint32 `json:"id"`
} //@name TransactionResponse

// swagger:parameters get_transaction
type _ struct {
	// The ID of a transaction
	// in:path
	ID string `json:"id"`
}

// TransactionsRequest Model swagger:parameters get_transactions
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

// TransactionItemsRequest Model swagger:parameters get_transactions
type TransactionItemsRequest struct {
	// Page
	// in:query
	Page uint32 `json:"page"`
	// Limit (max 100)
	// in:query
	Limit uint32 `json:"limit"`
	// in:query
	SortProperty string `json:"sortProperty"`
	// in:query
	SortDirection string `json:"sortDirection"`
}

// Transactions Model
type Transactions struct {
	Total int32         `json:"total" format:"int32"`
	Page  int32         `json:"page"  format:"int32"`
	Limit int32         `json:"limit" format:"int32"`
	Data  []Transaction `json:"data"`
} // @name Transactions

// TransactionItems Model
type TransactionItems struct {
	Total int32             `json:"total" format:"int32"`
	Page  int32             `json:"page"  format:"int32"`
	Limit int32             `json:"limit" format:"int32"`
	Data  []TransactionItem `json:"data"`
} // @name TransactionItems

// CustomersRequest Model swagger:parameters get_customers
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

// Customers Model
type Customers struct {
	Total int32      `json:"total" format:"int32"`
	Page  int32      `json:"page"  format:"int32"`
	Limit int32      `json:"limit" format:"int32"`
	Data  []Customer `json:"data"`
} // @name Customers

// ErrorResponse Model
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
} // @name ErrorResponse

// ErrorResponse400 Model
type ErrorResponse400 struct {
	Errors []ErrorMsg `json:"errors"`
} // @name ErrorResponse400

// Customer Model
type Customer struct {
	Model
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

// UpdateCustomer Model swagger:parameters update_customer
type UpdateCustomer struct {
	ID    uint32 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Role Model
type Role struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

// Address Model
type Address struct {
	Model
	Street     string `json:"street"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
	UserID     uint32 `json:"userId"`
}

// Business Model.
type Business struct {
	Model
	Name       string `json:"name"`
	HourlyRate uint32 `gorm:"column:hourly_rate" json:"hourlyRate"`
	Vat        uint32 `json:"vat"`
	Street     string `json:"street"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
}

// AccountData Model
type AccountData struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// User Model
type User struct {
	Model
	Username     string        `json:"username"`
	Address      Address       `json:"address"`
	Transactions []Transaction `json:"transactions"`
	Roles        []Role        `json:"roles"`
	BusinessID   uint32        `json:"businessId"`
	Account      AccountData   `json:"account"`
}

// UpdateAddress Model
type UpdateAddress struct {
	ID            uint32 `json:"id"`
	Street        string `json:"street"`
	PostalCode    string `json:"postalCode"`
	City          string `json:"city"`
	Country       string `json:"country"`
	LandLinePhone string `json:"landlinePhone"`
	MobilePhone   string `json:"mobilePhone"`
}

// UpdateUser Model swagger:parameters update_user
type UpdateUser struct {
	ID       uint32        `json:"id"`
	UserName string        `json:"username"`
	Address  UpdateAddress `json:"address"`
}

// UpdateBusiness Model
type UpdateBusiness struct {
	Name       string `json:"name"`
	HourlyRate uint32 `json:"hourlyRate"`
	Vat        uint32 `json:"vat"`
	Street     string `json:"street"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
}

// Currency Model
type Currency struct {
	Base 		string   `json:"base"`
	Symbol  string   `json:"symbol"`
}

// CurrencyResponse Model
type CurrencyResponse struct {
	To 		   string   	 `json:"to"`
	Base     string      `json:"base"`
	Value    string      `json:"value"`
}

// Lang Model
type Lang struct {
	ID 		   		uint32   	  `json:"id"`
	Language    string      `json:"language"`
	Currency    string      `json:"currency"`
} // @name Lang


// Langs Model
type Langs struct {
	Data  []Lang `json:"data"`
} // @name Langs
