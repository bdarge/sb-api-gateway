package lang

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bdarge/api-gateway/out/lang"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetLang godoc
// @Summary Get list of supported languages
// @ID lang
// @Success 200 {object} models.Langs
// @Router /lang [get]
func GetLang(ctx *gin.Context, client lang.LangServiceClient) {
	log.Printf("request uri %s", ctx.Request.RequestURI)

	res, err := client.GetLang(context.Background(), &lang.LangGetRequest{})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	message, err := protojson.Marshal(res)
	if err != nil {
		log.Printf("failed to cast type to bytes %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	log.Printf("message: %s", message)

	var data models.Langs
	err = json.Unmarshal(message, &data)
	if err != nil {
		log.Printf("failed to cast type, %v, %v", err, string(message))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
