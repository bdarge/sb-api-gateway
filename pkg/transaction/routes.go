package transaction

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

	svcItem := &TranItemServiceClient{
		Client: InitTranItemServiceClient(c),
	}

	routes := router.Group("/transaction")
	{
		routes.Use(a.AuthRequired)
		routes.POST("", svc.CreateTransaction)
		routes.GET("/:id", svc.GetTransaction)
		routes.GET("", svc.GetTransactions)
		routes.PATCH("/:id", svc.UpdateTransaction)
		routes.DELETE("/:id", svc.DeleteTransaction)

		routes.POST("/:id/item", svcItem.CreateTransactionItem)
		routes.GET("/:id/item", svcItem.GetTransactionItems)
		routes.GET("/:id/item/:item-id", svcItem.GetTransactionItem)
		routes.PATCH("/:id/item/:item-id", svcItem.UpdateTransactionItem)
		routes.DELETE("/:id/item/:item-id", svcItem.DeleteTransactionItem)
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

// CreateTransactionItem creates a pb
func (svc *TranItemServiceClient) CreateTransactionItem(ctx *gin.Context) {
	CreateTransactionItem(ctx, svc.Client)
}

// GetTransactionItem gets a pb
func (svc *TranItemServiceClient) GetTransactionItem(ctx *gin.Context) {
	GetTransactionItem(ctx, svc.Client)
}

// UpdateTransactionItem patch a transaction
func (svc *TranItemServiceClient) UpdateTransactionItem(ctx *gin.Context) {
	UpdateTransactionItem(ctx, svc.Client)
}

// DeleteTransactionItem delete a transaction item
func (svc *TranItemServiceClient) DeleteTransactionItem(ctx *gin.Context) {
	DeleteTransactionItem(ctx, svc.Client)
}

// GetTransactionItems all transaction items
func (svc *TranItemServiceClient) GetTransactionItems(ctx *gin.Context) {
	GetTransactionItems(ctx, svc.Client)
}
