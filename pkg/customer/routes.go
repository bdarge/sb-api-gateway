package customer

import (
	"github.com/bdarge/api-gateway/pkg/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register routes
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

// CreateCustomer create a customer
func (svc *ServiceClient) CreateCustomer(ctx *gin.Context) {
	CreateCustomer(ctx, svc.Client)
}

// GetCustomer returns a customer
func (svc *ServiceClient) GetCustomer(ctx *gin.Context) {
	GetCustomer(ctx, svc.Client)
}

// GetCustomers returns list of customers
func (svc *ServiceClient) GetCustomers(ctx *gin.Context) {
	GetCustomers(ctx, svc.Client)
}

// UpdateCustomer udpate a customer
func (svc *ServiceClient) UpdateCustomer(ctx *gin.Context) {
	UpdateCustomer(ctx, svc.Client)
}

// DeleteCustomer delete a customer
func (svc *ServiceClient) DeleteCustomer(ctx *gin.Context) {
	DeleteCustomer(ctx, svc.Client)
}
