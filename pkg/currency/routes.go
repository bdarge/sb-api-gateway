package currency

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

	routes := router.Group("/currency")
	{
		routes.Use(a.AuthRequired)
		routes.POST("", svc.Convert)
	}
}

// Convert returns currency rate
func (svc *ServiceClient) Convert(ctx *gin.Context) {
	Convert(ctx, svc.Client)
}

