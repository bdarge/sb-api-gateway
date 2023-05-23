package auth

import (
	"context"
	. "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Register
// @Summary Register a user
// @ID register
// @Param register body models.Account true "Add account details"
// @Success 200 {object} models.Account
// @Router /auth/register [post]
func Register(ctx *gin.Context, c AuthServiceClient) {
	body := models.Account{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.Register(context.Background(), &RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		if res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status), res.Error)
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{
					"error":   "ACTIONERR-1",
					"message": "An error happened, please check later."})
		}
		return
	}

	ctx.JSON(int(res.Status), &res)
}

// Login
// @Summary Authenticate a user
// @ID login
// @Param login body models.Login true "Add login credentials"
// @Success 200 {object} models.Login
// @Router /auth/login [post]
// @Security ApiKeyAuth
func Login(ctx *gin.Context, c AuthServiceClient) {
	b := models.Account{}
	log.Printf("Authenticate %v", ctx.Request.Body)
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	log.Printf("Request mapped, %v", b)

	res, err := c.Login(context.Background(), &LoginRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil || res.Status >= 400 {
		if res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status), res.Error)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
