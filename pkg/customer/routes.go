package customer

import (
	"github.com/bdarge/api-gateway/pkg/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := router.Group("/customer")
	{
		routes.Use(a.AuthRequired)
		routes.POST("/", svc.CreateCustomer)
		routes.GET("/:id", svc.GetCustomer)
	}
}

func (svc *ServiceClient) CreateCustomer(ctx *gin.Context) {
	CreateCustomer(ctx, svc.Client)
}

func (svc *ServiceClient) GetCustomer(ctx *gin.Context) {
	GetCustomer(ctx, svc.Client)
}
