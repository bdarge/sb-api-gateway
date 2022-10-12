package models

import (
	_ "github.com/bdarge/sb-api-gateway/cmd/docs"
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
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Order struct {
	Currency     string `json:"currency"`
	Description  string `json:"description"`
	DeliveryDate string `json:"deliveryDate"`
	CustomerId   int64  `json:"customerId"`
}
