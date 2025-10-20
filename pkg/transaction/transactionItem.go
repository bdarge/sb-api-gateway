package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
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
	utils.Logger()
	// read transaction id
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	item := models.TransactionItem{}

	if err := ctx.BindJSON(&item); err != nil {
		slog.Error("Failed to create a transaction", "error", err)
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

	slog.Info("Create transaction item", "item", item)

	inBytes, err := json.Marshal(item)
	if err != nil {
		slog.Error("Failed to create a transaction", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var data transactionItem.CreateTransactionItemRequest
	slog.Debug("stringfly data in bytes", "InBytes", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &data)

	if err != nil {
		slog.Error("Failed to create a transaction", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	slog.Debug("Response message", "Response", &data)

	res, err := c.CreateTransactionItem(context.Background(), &data)

	if res != nil && res.Status >= 400 {
		slog.Error("Failed to create a transaction", "error", err)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	if err != nil {
		slog.Error("Failed to create a transaction", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
	}

	ctx.JSON(http.StatusCreated, &models.CreateResponse{ID: res.Id})
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
	slog.Info("Delete transaction", "ID", id)

	res, err := c.DeleteTransactionItem(context.Background(), &transactionItem.DeleteTransactionItemRequest{
		Id: uint32(id),
	})

	if res != nil && res.Status >= 400 {
		slog.Error("Failed to delete", "error", res.Error)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	if err != nil {
		slog.Error("Failed to delete", "error", res.Error)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
	}

	ctx.Status(http.StatusNoContent)
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
	utils.Logger()
	slog.Info("request uri", "uri", ctx.Request.RequestURI)
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
		slog.Error("Failed to get transaction items", "error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	inBytes, err := json.Marshal(request)
	if err != nil {
		slog.Error("Failed to get transaction items", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	slog.Info("Stringfly data in bytes", "Bytes", inBytes)

	var requestMessage transactionItem.GetTransactionItemsRequest

	// ignore unknown fields
	// unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = protojson.Unmarshal(inBytes, &requestMessage)
	if err != nil {
		slog.Error("Failed to unmarshal to a request object", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	// set TransactionId
	requestMessage.TransactionId = uint32(id)

	slog.Debug("request message", "RequestMessage", &requestMessage)
	res, err := c.GetTransactionItems(context.Background(), &requestMessage)

	if res != nil && res.Status >= 400 {
		errorHeader := "ACTIONERR-1"
		if err != nil {
			slog.Error("Failed to get transaction items", "error", err)
		} else {
			errorHeader = "ACTIONERR-2"
			slog.Error("Failed to get transaction items", "error", res)
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   errorHeader,
				"message": "An error happened, please check later."})
		return
	}

	message, err := protojson.Marshal(res)

	if err != nil {
		slog.Error("Failed to cast response to bytes", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	slog.Info("Marshal to TransactionItems", "Response", message)
	var data models.TransactionItems
	err = json.Unmarshal(message, &data)

	if err != nil {
		slog.Error("Failed to cast to TransactionItems", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})

		return
	}

	ctx.JSON(http.StatusOK, data)
}

func toTransactionItemModel(data *model.TransactionItem) (*models.TransactionItem, error) {
	slog.Info("Marshal data", "data", data)
	message, err := protojson.Marshal(data)
	if err != nil {
		slog.Error("Failed to marsha data", "error", err)
		return nil, err
	}

	var d models.TransactionItem
	err = json.Unmarshal(message, &d)

	if err != nil {
		slog.Error("Failed to unmarshal to TransactionItem", "message", message, "error", err)
		return nil, err
	}

	return &d, nil
}
