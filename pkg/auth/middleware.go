package auth

import (
	"context"
	auth "github.com/bdarge/api-gateway/out/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"strings"
)

// Middleware struct to manage service client
type Middleware struct {
	svc *ServiceClient
}

// InitAuthMiddleware initialize Middleware
func InitAuthMiddleware(svc *ServiceClient) Middleware {
	return Middleware{svc}
}

// AuthRequired check auth requirement
func (c *Middleware) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		slog.Info("token not found")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := c.svc.Client.ValidateToken(context.Background(), &auth.ValidateTokenRequest{
		Token: token[1],
	})

	if err != nil || res.Status != http.StatusOK {
		slog.Error("token is invalid", "error", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}
