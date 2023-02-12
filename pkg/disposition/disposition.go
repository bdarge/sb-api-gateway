package disposition

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/bdarge/api-gateway/out/disposition"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
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
func CreateDisposition(ctx *gin.Context, c DispositionServiceClient) {
	disposition := models.Disposition{}

	if err := ctx.BindJSON(&disposition); err != nil {
		log.Printf("Error: %s", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) { /**/
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: models.GetErrorMsg(fe)}
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

	log.Printf("save disposition %v", disposition)

	inBytes, err := json.Marshal(disposition)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var data CreateDispositionRequest
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

	res, err := c.CreateDisposition(context.Background(), &data)

	if err != nil && res.Status >= 400 {
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
func GetDisposition(ctx *gin.Context, c DispositionServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetDisposition(context.Background(), &GetDispositionRequest{
		Id: uint32(id),
	})

	log.Printf("backend returned data: %v", res)

	if err != nil || res.Status >= 400 {
		if res.Status >= 400 {
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

	message, err := protojson.Marshal(res.Data)
	log.Printf("message %s", message)
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
func GetDispositions(ctx *gin.Context, c DispositionServiceClient) {
	var request = &models.DispositionsRequest{}

	err := ctx.ShouldBindQuery(&request)
	if err != nil {
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

	var requestMessage GetDispositionsRequest

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
	log.Printf("request message: %v", &requestMessage)
	res, err := c.GetDispositions(context.Background(), &requestMessage)

	if err != nil || res.Status >= 400 {
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Error: %v", res)
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
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
