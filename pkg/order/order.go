package order

import (
	"context"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/bdarge/sb-api-gateway/pkg/order/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateOrder
// @Summary Create an order
// @ID create_order
// @Param order body models.Order true "Order"
// @Success 200 {object} pb.CreateOrderResponse
// @Router /order [post]
// @Security ApiKeyAuth
func CreateOrder(ctx *gin.Context, c pb.OrderServiceClient) {
	body := models.Order{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, _ := ctx.Get("userId")

	res, err := c.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		CustomerId:   body.CustomerId,
		Description:  body.Description,
		DeliveryDate: body.DeliveryDate,
		CreatedBy:    userId.(int64),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func GetOrder(ctx *gin.Context, c pb.OrderServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetOrder(context.Background(), &pb.GetOrderRequest{
		Id: int64(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
