package disposition

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

	routes := router.Group("/disposition")
	{
		routes.Use(a.AuthRequired)
		routes.POST("/", svc.CreateDisposition)
		routes.GET("/", svc.GetDisposition)
	}
}

// CreateDisposition creates an disposition
func (svc *ServiceClient) CreateDisposition(ctx *gin.Context) {
	CreateDisposition(ctx, svc.Client)
}

// GetDisposition gets an disposition
func (svc *ServiceClient) GetDisposition(ctx *gin.Context) {
	GetDisposition(ctx, svc.Client)
}
