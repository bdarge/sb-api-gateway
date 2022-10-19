package request

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

	routes := router.Group("/request")
	{
		routes.Use(a.AuthRequired)
		routes.POST("/", svc.CreateRequest)
		routes.GET("/", svc.GetRequest)
	}
}

// CreateRequest creates an request
func (svc *ServiceClient) CreateRequest(ctx *gin.Context) {
	CreateRequest(ctx, svc.Client)
}

// GetRequest gets an request
func (svc *ServiceClient) GetRequest(ctx *gin.Context) {
	GetRequest(ctx, svc.Client)
}
