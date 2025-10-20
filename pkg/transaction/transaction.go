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
	"github.com/bdarge/api-gateway/out/transaction"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
)


// CreateTransaction create a transaction
// @Summary Create a transaction, an order or a quote
// @ID create_transaction
// @Param transaction body models.NewTransaction true "Add transaction"
// @Success 201 {object} models.CreateResponse
// @Router /transaction [post]
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func CreateTransaction(ctx *gin.Context, c transaction.TransactionServiceClient) {
	utils.Logger()
	t := models.NewTransaction{}

	if err := ctx.BindJSON(&t); err != nil {
		slog.Error("Failed to bind to models.NewTransaction", "error", err.Error())
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

	slog.Info("Create a new transaction", "NewTransaction", t)

	inBytes, err := json.Marshal(t)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var data transaction.CreateTransactionRequest
	slog.Debug("stringfly data in bytes", "InBytes", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &data)

	if err != nil {
		slog.Error("Failed to create request object", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	slog.Debug("Request object created", "Request", &data)

	res, err := c.CreateTransaction(context.Background(), &data)

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			slog.Error("Server Error", "error", res.Error)
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			log.Printf("Server Error: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	ctx.JSON(http.StatusCreated, &models.CreateResponse{ID: res.Id})
}

// GetTransaction get a transaction
// @Summary Get transaction
// @ID get_transaction
// @Success 200 {object} models.Transaction
// @Router /transaction/{id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func GetTransaction(ctx *gin.Context, c transaction.TransactionServiceClient) {
	utils.Logger()
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetTransaction(context.Background(), &transaction.GetTransactionRequest{
		Id: uint32(id),
	})

	slog.Info("Backend returned data", "TransactionResponse", res)

	if res != nil && res.Status >= 400 {
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	response, err := convertToModel(res.Data)

	if res != nil && res.Status >= 400 {
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
		return
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTransactions return transactions
// @Summary Get transactions
// @ID get_transactions
// @Param page query int false "Page"
// @Param limit query int false "Limit (max 100)"
// @Param requestType query string false "pass nothing, 'order' or 'quote'"
// @Success 200 {object} models.Transactions
// @Router /transaction [Get]
// @Security Bearer
func GetTransactions(ctx *gin.Context, c transaction.TransactionServiceClient) {
	utils.Logger()
	slog.Info("request", "url", ctx.Request.RequestURI)
	var request = &models.TransactionsRequest{}

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
			slog.Error("Failed to get transactions", "error", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{
					"error":   "ACTIONERR-1",
					"message": "An error happened, please check later."})
		}
		return
	}

	inBytes, err := json.Marshal(request)
	if err != nil {
		slog.Error("Failed to get transactions", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	slog.Debug("stringfly request in bytes", "bytes", inBytes)

	var requestMessage transaction.GetTransactionsRequest

	// ignore unknown fields
	// unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = protojson.Unmarshal(inBytes, &requestMessage)
	if err != nil {
		slog.Error("Failed to get transactions", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	slog.Info("Request", "TransactionsRequest", &requestMessage)
	res, err := c.GetTransactions(context.Background(), &requestMessage)

	if res != nil && res.Status >= 400 {
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: res.Error})
				return
	}
	if err != nil {
		slog.Error("Failed to get transactions", "error", err)

		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	
	message, err := protojson.Marshal(res)
	if err != nil {
		slog.Error("Failed to cast type to bytes", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	slog.Debug("Get transactions response", "response", message)

	var data models.Transactions
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

// UpdateTransaction updates a transaction
// @Summary Update a transaction
// @ID update_transaction
// @Param transaction body models.UpdateTransaction true "Update transaction"
// @Success 200 {object} models.Transaction
// @Router /transaction/{id} [Patch]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func UpdateTransaction(ctx *gin.Context, c transaction.TransactionServiceClient) {
	utils.Logger()
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	slog.Info("Update transaction", "ID", id)

	u := models.UpdateTransaction{}

	if err := ctx.BindJSON(&u); err != nil {
		slog.Error("Failed to update", "error", err)
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

	inBytes, err := json.Marshal(u)
	if err != nil {
		slog.Error("Failed to update", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var update transaction.UpdateTransactionData
	slog.Debug("stringfly data in bytes", "bytes", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &update)

	if err != nil {
		slog.Error("Failed to update", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	slog.Debug("Update request object", "Request", &update)

	res, err := c.UpdateTransaction(context.Background(), &transaction.UpdateTransactionRequest{
		Id:   uint32(id),
		Data: &update,
	})

	response, err := convertToModel(res.Data)

	if res != nil && res.Status >= 400 {
		slog.Error("Failed to update", "error", err)
		ctx.AbortWithStatusJSON(int(res.Status),
			models.ErrorResponse{
				Error:   "ACTIONERR-2",
				Message: "An error happened, please check later."})
		return
	}

	if err != nil {
		slog.Error("Failed to update", "error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteTransaction delete
// @Summary Delete a transaction
// @ID delete_transaction
// @Success 200
// @Router /transaction/{id} [Delete]
// @Failure 500 {object} ErrorResponse
// @Security Bearer
func DeleteTransaction(ctx *gin.Context, c transaction.TransactionServiceClient) {
	utils.Logger()
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	slog.Info("Delete transaction", "ID", id)

	res, err := c.DeleteTransaction(context.Background(), &transaction.DeleteTransactionRequest{
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

func convertToModel(data *model.TransactionData) (*models.Transaction, error) {
	message, err := protojson.Marshal(data)
	slog.Info("Unmarshal", "Transaction", string(message))

	var d models.Transaction
	err = json.Unmarshal(message, &d)

	if err != nil {
		slog.Error("Failed to cast type", "data", string(message), "error", err)
		return nil, err
	}

	slog.Info("Unmarshaled", "Transaction", d)
	return &d, nil
}
