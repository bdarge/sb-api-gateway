package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/bdarge/api-gateway/out/model"
	"github.com/bdarge/api-gateway/out/transactionItem"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateTransactionItem godoc
// @Summary Create a transaction, an order or a quote
// @ID create_transaction_item
// @Param transactionItem body models.NewTransactionItem true "Add transaction item"
// @Success 201 {object} models.CreateResponse
// @Router /transaction/{id}/item [post]
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func CreateTransactionItem(ctx *gin.Context, c transactionItem.TransactionItemServiceClient) {

	// read transaction id
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	item := models.TransactionItem{}

	if err := ctx.BindJSON(&item); err != nil {
		log.Printf("Error: %s", err)
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
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

	// set the transaction id
	item.TransactionID = uint32(id)

	log.Printf("Create transaction item %v", item)

	inBytes, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var data transactionItem.CreateTransactionItemRequest
	log.Printf("stringfly data in bytes: %s", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &data)

	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("mesage: %v", &data)

	res, err := c.CreateTransactionItem(context.Background(), &data)

	if err == nil && res != nil && res.Status < 400 {
		ctx.JSON(http.StatusCreated, &models.CreateResponse{ID: res.Id})
		return
	}

	if res != nil && res.Status >= 400 {
		log.Printf("Server Error: %v", res.Error)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	log.Printf("Server Error: %v", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError,
		models.ErrorResponse{
			Error:   "ACTIONERR-1",
			Message: "An error happened, please check later."})
}

// GetTransactionItem godoc
// @Summary Get transaction item
// @ID get_transaction_item
// @Success 200 {object} models.TransactionItem
// @Router /transaction/{id}/item/{item-id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func GetTransactionItem(ctx *gin.Context, c transactionItem.TransactionItemServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("item-id"), 10, 32)

	if res, err := c.GetTransactionItem(context.Background(), &transactionItem.GetTransactionItemRequest{
		Id: uint32(id),
	}); err == nil {

		log.Printf("backend returned data: %v", res)
		if response, err := toTransactionItemModel(res.Data); err == nil && res != nil && res.Status < 400 {
			ctx.JSON(http.StatusOK, response)
			return
		}

		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError,
		models.ErrorResponse{
			Error:   "ACTIONERR-1",
			Message: "An error happened, please check later."})
}

// UpdateTransactionItem godoc
// @Summary Update a transaction item
// @ID update_transaction_item
// @Param transaction body models.UpdateTransactionItem true "Update transaction item"
// @Success 200 {object} models.Transaction
// @Router /transaction/{id}/item/{item-id} [Patch]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func UpdateTransactionItem(ctx *gin.Context, c transactionItem.TransactionItemServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("item-id"), 10, 32)
	log.Printf("Update transaction item (id = %d)", id)
	u := models.TransactionItem{}

	if err := ctx.BindJSON(&u); err != nil {
		log.Printf("Failed to bind request to json object: %s", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "ACTIONERR-3",
			"message": "One or more fileds have invalid data type",
		})
		return
	}

	inBytes, err := json.Marshal(u)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var update transactionItem.UpdateTransactionItem
	log.Printf("stringfly data in bytes: %s", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &update)

	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("mesage: %v", &update)

	res, err := c.UpdateTransactionItem(context.Background(), &transactionItem.UpdateTransactionItemRequest{
		Id:   uint32(id),
		Data: &update,
	})

	if err != nil {
		log.Printf("Error when updating: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	response, err := toTransactionItemModel(res.Data)

	if err == nil && res != nil && res.Status < 400 {
		ctx.JSON(http.StatusOK, response)
		return
	}

	if res != nil && res.Status >= 400 {
		log.Printf("Server Error: %v", res.Error)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: "An error happened, please check later."})
		return
	}

	log.Printf("Error while mapping data: %v", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError,
		models.ErrorResponse{
			Error:   "ACTIONERR-1",
			Message: "An error happened, please check later."})
}

// DeleteTransactionItem delete
// @Summary Delete a transaction item
// @ID delete_transaction_item
// @Success 200
// @Router /transaction/{id}/item/{item-id} [Delete]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func DeleteTransactionItem(ctx *gin.Context, c transactionItem.TransactionItemServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("item-id"), 10, 32)
	log.Printf("Delete transaction with id: %d", id)

	res, err := c.DeleteTransactionItem(context.Background(), &transactionItem.DeleteTransactionItemRequest{
		Id: uint32(id),
	})

	if err == nil && res != nil && res.Status < 400 {
		ctx.Status(http.StatusNoContent)
		return
	}

	if res != nil && res.Status >= 400 {
		log.Printf("Server Error: %v", res.Error)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	log.Printf("Server Error: %v", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError,
		models.ErrorResponse{
			Error:   "ACTIONERR-1",
			Message: "An error happened, please check later."})
}

// GetTransactionItems return transaction items
// @Summary Get transactions
// @ID get_transactions
// @Param page query int false "Page"
// @Param limit query int false "Limit (max 100)"
// @Param requestType query string false "pass nothing, 'order' or 'quote'"
// @Success 200 {object} models.Transactions
// @Router /transaction [Get]
// @Security Bearer
func GetTransactionItems(ctx *gin.Context, c transactionItem.TransactionItemServiceClient) {
	log.Printf("request uri %s", ctx.Request.RequestURI)
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)
	var request = &models.TransactionItemsRequest{}

	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
			return
		}
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
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
	log.Printf("stringfly data in bytes: %s", inBytes)

	var requestMessage transactionItem.GetTransactionItemsRequest

	// ignore unknown fields
	// unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = protojson.Unmarshal(inBytes, &requestMessage)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	// set TransactionId
	requestMessage.TransactionId = uint32(id)

	log.Printf("request message: %v", &requestMessage)
	res, err := c.GetTransactionItems(context.Background(), &requestMessage)

	if err != nil || (res != nil && res.Status >= 400) {
		errorHeader := "ACTIONERR-1"
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			errorHeader = "ACTIONERR-2"
			log.Printf("Error: %v", res)
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   errorHeader,
				"message": "An error happened, please check later."})
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
	var data models.TransactionItems
	err = json.Unmarshal(message, &data)

	if err == nil {
		ctx.JSON(http.StatusOK, data)
		return
	}

	log.Printf("failed to cast type, %v: %v", string(message), err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{
			"error":   "ACTIONERR-1",
			"message": "An error happened, please check later."})
	return
}

func toTransactionItemModel(data *model.TransactionItem) (*models.TransactionItem, error) {
	message, err := protojson.Marshal(data)
	log.Printf("message %s", message)

	if err != nil {
		log.Printf("failed to marsha data, %v: %v", data, err)
		return nil, err
	}

	var d models.TransactionItem
	err = json.Unmarshal(message, &d)

	if err != nil {
		log.Printf("failed to cast type, %v: %v", string(message), err)
		return nil, err
	}

	return &d, nil
}
