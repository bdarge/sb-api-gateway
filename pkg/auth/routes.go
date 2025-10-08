package auth

import (
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register routes
func RegisterRoutes(r *gin.RouterGroup, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
		Config: *c,
	}

	routes := r.Group("/auth")
	routes.POST("/register", svc.Register)
	routes.POST("/login", svc.Login)
	routes.POST("/refresh-token", svc.refreshToken)
	return svc
}

// Register a new user
func (svc *ServiceClient) Register(ctx *gin.Context) {
	Register(ctx, svc.Client)
}

// Login a user
func (svc *ServiceClient) Login(ctx *gin.Context) {
	Login(ctx, svc.Client, &svc.Config)
}

// Refresh token
func (svc *ServiceClient) refreshToken(ctx *gin.Context) {
	RefreshToken(ctx, svc.Client)
}
