package profile

import (
	"context"
	"encoding/json"
	. "github.com/bdarge/api-gateway/out/profile"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"strconv"
)

// GetUser
// @Summary Get User
// @ID get_user
// @Success 200 {object} models.User
// @Router /user/{id} [Get]
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func GetUser(ctx *gin.Context, c ProfileServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetUser(context.Background(), &GetUserRequest{
		Id: uint32(id),
	})

	log.Printf("backend returned data: %v", res)

	if err != nil || res.Status >= 400 {
		if res.Status >= 400 {
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

	response, err := convertToModel(res.Data)

	if err != nil || res.Status >= 400 {
		if res.Status >= 400 {
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

	ctx.JSON(http.StatusOK, response)
}

func convertToModel(data *UserData) (*models.User, error) {
	message, err := protojson.Marshal(data)
	log.Printf("message %s", message)

	var d models.User
	err = json.Unmarshal(message, &d)

	if err != nil {
		log.Printf("failed to cast type, %v, %v", err, string(message))
		return nil, err
	}

	return &d, nil
}
