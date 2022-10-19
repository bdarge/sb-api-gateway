package request

import (
	"context"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/bdarge/sb-api-gateway/pkg/request/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateRequest
// @Summary Create an request
// @ID create_request
// @Param request body models.request true "Add request"
// @Success 201 {object} pb.CreateRequestResponse
// @Router /request [post]
// @Security ApiKeyAuth
func CreateRequest(ctx *gin.Context, c pb.RequestServiceClient) {
	order := models.Request{}

	if err := ctx.BindJSON(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.CreateRequest(context.Background(), &pb.CreateRequestRequest{
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

// GetRequest
// @Summary Get a request
// @ID create_request
// @Param request body models.request true "Add request"
// @Success 200 {object} pb.GetRequestResponse
// @Router /request [post]
// @Security ApiKeyAuth
func GetRequest(ctx *gin.Context, c pb.RequestServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetRequest(context.Background(), &pb.GetRequestRequest{
		Id: int64(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
