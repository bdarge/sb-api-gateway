package customer

import (
	"context"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/customer/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
func CreateCustomer(ctx *gin.Context, client pb.CustomerServiceClient) {
	disposition := models.Customer{}

	if err := ctx.BindJSON(&disposition); err != nil {
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

	res, err := client.CreateCustomer(context.Background(), &pb.CreateCustomerRequest{
		Name:  disposition.Name,
		Email: disposition.Email,
	})

	if err != nil {
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
func GetCustomer(ctx *gin.Context, client pb.CustomerServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := client.GetCustomer(context.Background(), &pb.GetCustomerRequest{
		Id: id,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, res.Data)
}
