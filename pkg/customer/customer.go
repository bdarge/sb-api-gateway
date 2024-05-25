package customer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	. "github.com/bdarge/api-gateway/out/customer"
	"github.com/bdarge/api-gateway/out/model"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateCustomer create a customer
// @Summary Create a customer
// @ID create_customer
// @Param customer body models.Customer true "Add customer"
// @Success 201 {object} models.CreateResponse
// @Router /customer [post]
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func CreateCustomer(ctx *gin.Context, client CustomerServiceClient) {
	customer := models.Customer{}

	if err := ctx.BindJSON(&customer); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		}
		return
	}

	res, err := client.CreateCustomer(context.Background(), &CreateCustomerRequest{
		Name:  customer.Name,
		Email: customer.Email,
	})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	ctx.JSON(http.StatusCreated, &models.CreateResponse{ID: res.Id})
}

// GetCustomer godoc
// @Summary Get customer
// @ID get_customer
// @Success 200 {object} models.Customer
// @Router /customer/{id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func GetCustomer(ctx *gin.Context, client CustomerServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := client.GetCustomer(context.Background(), &GetCustomerRequest{
		Id: uint32(id),
	})

	if err != nil {
		log.Printf("error getting a customer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, res.Data)
}

// GetCustomers godoc
// @Summary Get customers
// @ID get_customers
// @Param page query int false "Page"
// @Param limit query int false "Limit (max 100)"
// @Success 200 {object} models.Customers
// @Router /customer [Get]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func GetCustomers(ctx *gin.Context, client CustomerServiceClient) {
	log.Printf("request uri %s", ctx.Request.RequestURI)

	var request = &models.CustomersRequest{}

	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		} else {
			log.Printf("Error: %v", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{
					"error":   "ACTIONERR-1",
					"message": "An error happened, please check later."})
		}
		return
	}

	inBytes, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("stringify data in bytes: %s", inBytes)

	var requestMessage GetCustomersRequest

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &requestMessage)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("request message: %v", &requestMessage)
	res, err := client.GetCustomers(context.Background(), &requestMessage)

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	message, err := protojson.Marshal(res)
	if err != nil {
		log.Printf("failed to cast type to bytes %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	log.Printf("message: %s", message)

	var data models.Customers
	err = json.Unmarshal(message, &data)
	if err != nil {
		log.Printf("failed to cast type, %v, %v", err, string(message))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// UpdateCustomer godoc
// @Summary Update a customer
// @ID update_customer
// @Param customer body models.UpdateCustomer true "Update customer"
// @Success 200 {object} models.Customer
// @Router /customer/{id} [Patch]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func UpdateCustomer(ctx *gin.Context, c CustomerServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)
	u := models.UpdateCustomer{}

	if err := ctx.BindJSON(&u); err != nil {
		log.Printf("Error: %s", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) { /**/
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "ACTIONERR-1",
				"message": err.Error(),
			})
		}
		return
	}

	log.Printf("update customer")

	inBytes, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to marshal update data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var update UpdateCustomerData
	log.Printf("stringfly data in bytes: %s", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &update)

	if err != nil {
		log.Printf("Failed to unmarsha to proto type: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("mesage: %v", &update)

	res, err := c.UpdateCustomer(context.Background(), &UpdateCustomerRequest{
		Id:   uint32(id),
		Data: &update,
	})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			log.Printf("Failed to updated customer: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	response, err := convertToModel(res.Data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @ID delete_customer
// @Success 200
// @Router /Customer/{id} [Delete]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func DeleteCustomer(ctx *gin.Context, c CustomerServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.DeleteCustomer(context.Background(), &DeleteCustomerRequest{
		Id: uint32(id),
	})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
	return
}

func convertToModel(data *model.CustomerData) (*models.Customer, error) {
	message, err := protojson.Marshal(data)
	log.Printf("message %s", message)

	var d models.Customer
	err = json.Unmarshal(message, &d)

	if err != nil {
		log.Printf("failed to cast type, %v, %v", err, string(message))
		return nil, err
	}

	return &d, nil
}
