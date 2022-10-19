package disposition

import (
	"context"
	"github.com/bdarge/sb-api-gateway/pkg/disposition/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
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
	order := models.Disposition{}

	if err := ctx.BindJSON(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.CreateDisposition(context.Background(), &pb.CreateDispositionRequest{
		CustomerId:   order.CustomerId,
		Description:  order.Description,
		DeliveryDate: order.DeliveryDate,
		CreatedBy:    order.CreatedBy,
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
