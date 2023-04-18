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
	}
}

// GetUser gets a user
func (svc *ServiceClient) GetUser(ctx *gin.Context) {
	GetUser(ctx, svc.Client)
}
