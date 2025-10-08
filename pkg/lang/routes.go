package lang

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

	routes := router.Group("/lang")
	{
		routes.Use(a.AuthRequired)
		routes.GET("", svc.getLang)
	}
}

// GetLang returns list of supported languages
func (svc *ServiceClient) getLang(ctx *gin.Context) {
	GetLang(ctx, svc.Client)
}
