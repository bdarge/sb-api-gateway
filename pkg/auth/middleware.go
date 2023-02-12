package auth

import (
	"context"
	auth "github.com/bdarge/api-gateway/out/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Middleware struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) Middleware {
	return Middleware{svc}
}

func (c *Middleware) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := c.svc.Client.Validate(context.Background(), &auth.ValidateRequest{
		Token: token[1],
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}
