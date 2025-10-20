package auth

import (
	"context"
	"net/http"

	"github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

// Account to register
type Account struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
} // @name Account



// LogValue redact sensitive values
func (a Account) LogValue() slog.Value {
	d := slog.GroupValue(
		slog.String("email", "[redacted]"),
		slog.String("password", "[redacted]"),
	)
	return d
}


// Register godoc
// @Summary Register a user
// @ID register
// @Param register body Account true "Add account details"
// @Success 200 {object} Account
// @Router /auth/register [post]
func Register(ctx *gin.Context, c auth.AuthServiceClient) {
	body := Account{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	res, err := c.Register(context.Background(), &auth.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		if res != nil && res.Status >= 400 {
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
func Login(ctx *gin.Context, authClient auth.AuthServiceClient, config *config.Config) {
	utils.Logger()

	account := Account{}
	slog.Info("Start authenticating with email and password")
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": "Invalid inputs. Please check your inputs"})
		return
	}

	slog.Info("Request mapped", "account", account)

	result, err := authClient.Login(context.Background(), &auth.LoginRequest{
		Email:    account.Email,
		Password: account.Password,
	})

	if err != nil {
		if result != nil && result.Status >= 400 {
			ctx.AbortWithStatusJSON(int(result.Status), result.Error)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	res := &auth.LoginResponse{
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
func RefreshToken(ctx *gin.Context, c auth.AuthServiceClient) {
	utils.Logger()

	token, err := ctx.Cookie("token")

	if err != nil {
		slog.Error("Cookie not found", "error", err)
		ctx.AbortWithStatusJSON(403, models.ErrorResponse{
					Error:   "ACTIONERR-3",
					Message: "Not authorized"})
		return
	}
	res, err := c.RefreshToken(context.Background(), &auth.RefreshTokenRequest{
		Token: token,
	})

	if err != nil {
		slog.Error("Failed to refresh token", "error", err)
		ctx.AbortWithStatusJSON(403, models.ErrorResponse{
					Error:   "ACTIONERR-3",
					Message: "Not authorized"})
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
