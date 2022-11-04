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

	routes := router.Group("/pb")
	{
		routes.Use(a.AuthRequired)
		routes.POST("/", svc.CreateDisposition)
		routes.GET("/:id", svc.GetDisposition)
		routes.GET("/", svc.GetDispositions)
	}
}

// CreateDisposition creates a pb
func (svc *ServiceClient) CreateDisposition(ctx *gin.Context) {
	CreateDisposition(ctx, svc.Client)
}

// GetDisposition gets a pb
func (svc *ServiceClient) GetDisposition(ctx *gin.Context) {
	GetDisposition(ctx, svc.Client)
}

// GetDispositions all dispositions
func (svc *ServiceClient) GetDispositions(ctx *gin.Context) {
	GetDispositions(ctx, svc.Client)
}
