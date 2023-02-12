package auth

import (
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/auth")
	routes.POST("/register", svc.Register)
	routes.POST("/login", svc.Login)

	return svc
}

// Register a new user
func (svc *ServiceClient) Register(ctx *gin.Context) {
	Register(ctx, svc.Client)
}

// Login a user
func (svc *ServiceClient) Login(ctx *gin.Context) {
	Login(ctx, svc.Client)
}
