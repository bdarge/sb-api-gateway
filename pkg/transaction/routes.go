package transaction

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

	routes := router.Group("/transaction")
	{
		routes.Use(a.AuthRequired)
		routes.POST("", svc.CreateTransaction)
		routes.GET("/:id", svc.GetTransaction)
		routes.GET("", svc.GetTransactions)
		routes.PATCH("/:id", svc.UpdateTransaction)
		routes.DELETE("/:id", svc.DeleteTransaction)
	}
}

// CreateTransaction creates a pb
func (svc *ServiceClient) CreateTransaction(ctx *gin.Context) {
	CreateTransaction(ctx, svc.Client)
}

// GetTransaction gets a pb
func (svc *ServiceClient) GetTransaction(ctx *gin.Context) {
	GetTransaction(ctx, svc.Client)
}

// GetTransactions all transactions
func (svc *ServiceClient) GetTransactions(ctx *gin.Context) {
	GetTransactions(ctx, svc.Client)
}

// UpdateTransaction patch a transaction
func (svc *ServiceClient) UpdateTransaction(ctx *gin.Context) {
	UpdateTransaction(ctx, svc.Client)
}

// DeleteTransaction delete a transaction
func (svc *ServiceClient) DeleteTransaction(ctx *gin.Context) {
	DeleteTransaction(ctx, svc.Client)
}
