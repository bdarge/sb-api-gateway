package disposition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/disposition/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"strconv"
)

// CreateDisposition
// @Summary Create a disposition, an order or a quote
// @ID create_disposition
// @Param disposition body models.Disposition true "Add dispositions"
// @Success 201 {object} models.CreateResponse
// @Router /disposition [post]
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func CreateDisposition(ctx *gin.Context, c pb.DispositionServiceClient) {
	disposition := models.Disposition{}

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

	res, err := c.CreateDisposition(context.Background(), &pb.CreateDispositionRequest{
		CustomerId:   disposition.CustomerId,
		Description:  disposition.Description,
		DeliveryDate: timestamppb.New(disposition.DeliveryDate),
		CreatedBy:    disposition.CreatedBy,
		RequestType:  disposition.RequestType,
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

// GetDisposition
// @Summary Get disposition
// @ID get_disposition
// @Success 200 {object} models.Disposition
// @Router /disposition/{id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func GetDisposition(ctx *gin.Context, c pb.DispositionServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetDisposition(context.Background(), &pb.GetDispositionRequest{
		Id: id,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "ACTIONERR-1",
				Message: "An error happened, please check later."})
		return
	}

	message, err := protojson.Marshal(res.Data)
	var data models.Disposition
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

// GetDispositions
// @Summary Get dispositions
// @ID get_dispositions
// @Param page query int false "Page"
// @Param limit query int false "Limit (max 100)"
// @Param requestType query string false "pass nothing, 'order' or 'quote'"
// @Success 200 {object} models.Dispositions
// @Router /disposition [Get]
// @Security ApiKeyAuth
func GetDispositions(ctx *gin.Context, c pb.DispositionServiceClient) {
	var requestType = ctx.Param("requestTye")
	res, err := c.GetDispositions(context.Background(), &pb.GetDispositionsRequest{
		RequestType: requestType,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	// https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson
	message, err := protojson.Marshal(res)
	if err != nil {
		log.Printf("failed to cast type to bytes %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	var data models.Dispositions
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
