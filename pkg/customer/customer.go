package customer

import (
	"context"
	"errors"
	. "github.com/bdarge/api-gateway/out/customer"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

// CreateCustomer
// @Summary Create a customer
// @ID create_customer
// @Param customer body models.Customer true "Add customer"
// @Success 201 {object} models.CreateResponse
// @Router /customer [post]
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func CreateCustomer(ctx *gin.Context, client CustomerServiceClient) {
	customer := models.Customer{}

	if err := ctx.BindJSON(&customer); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: models.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		}
		return
	}

	res, err := client.CreateCustomer(context.Background(), &CreateCustomerRequest{
		Name:  customer.Name,
		Email: customer.Email,
	})

	if err != nil {
		log.Printf("error creating a cutomer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusCreated, &models.CreateResponse{ID: res.Id})
}

// GetCustomer
// @Summary Get customer
// @ID get_customer
// @Success 200 {object} models.Customer
// @Router /customer/{id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func GetCustomer(ctx *gin.Context, client CustomerServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := client.GetCustomer(context.Background(), &GetCustomerRequest{
		Id: uint32(id),
	})

	if err != nil {
		log.Printf("error getting a cutomer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, res.Data)
}
