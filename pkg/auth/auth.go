package auth

import (
	"context"
	"net/http"

	"github.com/bdarge/sb-api-gateway/pkg/auth/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
)

// Register
// @Summary Register a user
// @ID register
// @Param register body models.Account true "Add account details"
// @Success 200 {object} models.Account
// @Router /auth/register [post]
func Register(ctx *gin.Context, c pb.AuthServiceClient) {
	body := models.Account{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.Register(context.Background(), &pb.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-2",
				"message": "An error happened, please check later."})
		return
	}

	ctx.JSON(int(res.Status), &res)
}

// Login
// @Summary Authenticate a user
// @ID login
// @Param login body models.Login true "Add login credentials"
// @Success 200 {object} models.Account
// @Router /auth/login [post]
// @Security ApiKeyAuth
func Login(ctx *gin.Context, c pb.AuthServiceClient) {
	b := models.Account{}

	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.Login(context.Background(), &pb.LoginRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}