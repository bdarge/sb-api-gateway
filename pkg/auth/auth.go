package auth

import (
	"context"
	. "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
)

// Register godoc
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

// Login godoc
// @Summary Authenticate a user
// @ID login
// @Param login body models.Login true "Add login credentials"
// @Success 200 {object} models.LoginResponse
// @Router /auth/login [post]
func Login(ctx *gin.Context, authClient AuthServiceClient, config *config.Config) {
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

	result, err := authClient.Login(context.Background(), &LoginRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil || result.Status >= 400 {
		if err != nil && result.Status >= 400 {
			ctx.AbortWithStatusJSON(int(result.Status), result.Error)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	res := &LoginResponse{
		Status: result.Status,
		Token:  result.Token,
		Error:  result.Error,
	}
	ctx.SetSameSite(http.SameSiteNoneMode)

	ctx.SetCookie(
		"token", result.RefreshToken, config.RefreshTokenExpOn,
		"/", config.UIDomain, true, true,
	)

	ctx.JSON(http.StatusOK, &res)
}

// RefreshToken godoc
// @Summary Authenticate a user
// @ID refresh_token
// @Success 200 {object} models.LoginResponse
// @Router /auth/refresh-token [post]
func RefreshToken(ctx *gin.Context, c AuthServiceClient, _ *config.Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	token, err := ctx.Cookie("token")
	if err != nil {
		slog.Error("Cookie not found", "error", err)
		ctx.AbortWithStatusJSON(403, "Not authorized")
		return
	}
	res, err := c.RefreshToken(context.Background(), &RefreshTokenRequest{
		Token: token,
	})

	if err != nil || res.Status >= 400 {
		slog.Error("Failed to refresh token", "error", err)
		ctx.AbortWithStatusJSON(403, "Not authorized")
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
