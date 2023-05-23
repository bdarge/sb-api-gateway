package profile

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

	routes := router.Group("/user")
	{
		routes.Use(a.AuthRequired)
		routes.GET("/:id", svc.GetUser)
		routes.PATCH("/:id", svc.UpdateUser)
	}

	businessRoutes := router.Group("/business")
	{
		businessRoutes.Use(a.AuthRequired)
		businessRoutes.GET("/:id", svc.GetBusiness)
		businessRoutes.PATCH("/:id", svc.UpdateBusiness)
	}
}

// GetUser gets a user
func (svc *ServiceClient) GetUser(ctx *gin.Context) {
	GetUser(ctx, svc.Client)
}

// UpdateUser update a user
func (svc *ServiceClient) UpdateUser(ctx *gin.Context) {
	UpdateUser(ctx, svc.Client)
}

// GetBusiness gets a business
func (svc *ServiceClient) GetBusiness(ctx *gin.Context) {
	GetBusiness(ctx, svc.Client)
}

// UpdateBusiness updates a business
func (svc *ServiceClient) UpdateBusiness(ctx *gin.Context) {
	UpdateBusiness(ctx, svc.Client)
}
