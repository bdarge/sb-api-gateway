package order

import (
	"github.com/bdarge/sb-api-gateway/pkg/auth"
	"github.com/bdarge/sb-api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := router.Group("/order")
	{
		routes.Use(a.AuthRequired)
		routes.POST("/", svc.CreateOrder)
		routes.GET("/", svc.GetOrder)
	}
}

// CreateOrder creates an order
func (svc *ServiceClient) CreateOrder(ctx *gin.Context) {
	CreateOrder(ctx, svc.Client)
}

// GetOrder gets an order
func (svc *ServiceClient) GetOrder(ctx *gin.Context) {
	GetOrder(ctx, svc.Client)
}
