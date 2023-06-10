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
		routes.POST("", svc.CreateCustomer)
		routes.GET("/:id", svc.GetCustomer)
		routes.GET("", svc.GetCustomers)
		routes.PATCH("/:id", svc.UpdateCustomer)
		routes.DELETE("/:id", svc.DeleteCustomer)
	}
}

func (svc *ServiceClient) CreateCustomer(ctx *gin.Context) {
	CreateCustomer(ctx, svc.Client)
}

func (svc *ServiceClient) GetCustomer(ctx *gin.Context) {
	GetCustomer(ctx, svc.Client)
}

func (svc *ServiceClient) GetCustomers(ctx *gin.Context) {
	GetCustomers(ctx, svc.Client)
}

func (svc *ServiceClient) UpdateCustomer(ctx *gin.Context) {
	UpdateCustomer(ctx, svc.Client)
}

func (svc *ServiceClient) DeleteCustomer(ctx *gin.Context) {
	DeleteCustomer(ctx, svc.Client)
}
