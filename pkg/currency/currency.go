package currency

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/bdarge/api-gateway/out/currency"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Convert convert a currency
func Convert(ctx *gin.Context, client currency.CurrencyClient) {
	log.Printf("Convert %v", ctx.Request.Body)

	c := models.Currency{}

	if err := ctx.BindJSON(&c); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		}
		return
	}

	res, err := client.Convert(context.Background(), &currency.CurrencyRequest{
		Base:  c.Base,
		Symbol: c.Symbol,
	})

	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
