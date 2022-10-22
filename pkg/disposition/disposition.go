package disposition

import (
	"context"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/disposition/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

// CreateDisposition
// @Summary Create a disposition
// @ID create_request
// @Param disposition body models.disposition true "Add disposition"
// @Success 201 {object} pb.CreateDispositionResponse
// @Router /disposition [post]
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	res, err := c.CreateDisposition(context.Background(), &pb.CreateDispositionRequest{
		CustomerId:   disposition.CustomerId,
		Description:  disposition.Description,
		DeliveryDate: disposition.DeliveryDate,
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

	ctx.JSON(http.StatusCreated, &res)
}

// GetDisposition
// @Summary Get a disposition
// @ID create_request
// @Param disposition body models.disposition true "Add disposition"
// @Success 200 {object} pb.GetDispositionResponse
// @Router /disposition [post]
// @Security ApiKeyAuth
func GetDisposition(ctx *gin.Context, c pb.DispositionServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetDisposition(context.Background(), &pb.GetDispositionRequest{
		Id: int64(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
