package profile

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api-gateway/out/model"
	. "github.com/bdarge/api-gateway/out/profile"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// UpdateUser
// @Summary Update a user
// @ID update_user
// @Param user body models.UpdateUser true "Update user"
// @Success 200 {object} models.User
// @Router /user/{id} [Patch]
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
func UpdateUser(ctx *gin.Context, c ProfileServiceClient) {
	log.Printf("** update user %v", ctx.Request.FormValue("id"))
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)
	u := models.UpdateUser{}

	if err := ctx.BindJSON(&u); err != nil {
		log.Printf("Error: %s", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) { /**/
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: models.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "ACTIONERR-1",
				"message": err.Error(),
			})
		}
		return
	}

	inBytes, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to marshal update data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var updateUserData UpdateUserData
	log.Printf("stringfly data in bytes: %s", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &updateUserData)

	if err != nil {
		log.Printf("Failed to unmarsha to proto type: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	updateUserData.Id = uint32(id)
	log.Printf("mesage: %v", &updateUserData)

	res, err := c.UpdateUser(context.Background(), &UpdateUserRequest{
		Id:   uint32(id),
		Data: &updateUserData,
	})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			log.Printf("Failed to updated user: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	response, err := convertToModel(res.Data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func GetBusiness(ctx *gin.Context, c ProfileServiceClient) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := c.GetBusiness(context.Background(), &GetBusinessRequest{
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

	response, err := convertToBusinessModel(res.Data)

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

func UpdateBusiness(ctx *gin.Context, c ProfileServiceClient) {
	log.Printf("** update business %v", ctx.Request.FormValue("id"))
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)
	u := models.UpdateBusiness{}

	if err := ctx.BindJSON(&u); err != nil {
		log.Printf("Error: %s", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) { /**/
			out := make([]models.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = models.ErrorMsg{Field: fe.Field(), Message: models.GetErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse400{Errors: out})
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "ACTIONERR-1",
				"message": err.Error(),
			})
		}
		return
	}

	inBytes, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to marshal update data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	var updateBusinessData UpdateBusinessData
	log.Printf("stringfly data in bytes: %s", inBytes)

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(inBytes, &updateBusinessData)

	if err != nil {
		log.Printf("Failed to unmarsha to proto type: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}
	log.Printf("mesage: %v", &updateBusinessData)

	res, err := c.UpdateBusiness(context.Background(), &UpdateBusinessRequest{
		Id:   uint32(id),
		Data: &updateBusinessData,
	})

	if err != nil || res.Status >= 400 {
		if res != nil && res.Status >= 400 {
			ctx.AbortWithStatusJSON(int(res.Status),
				models.ErrorResponse{
					Error:   "ACTIONERR-2",
					Message: res.Error})
		} else {
			log.Printf("Failed to updated user: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{
					Error:   "ACTIONERR-1",
					Message: "An error happened, please check later."})
		}
		return
	}

	response, err := convertToBusinessModel(res.Data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error":   "ACTIONERR-1",
				"message": "An error happened, please check later."})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func convertToBusinessModel(data *model.BusinessData) (*models.Business, error) {
	message, err := protojson.Marshal(data)
	log.Printf("message %s", message)

	var d models.Business
	err = json.Unmarshal(message, &d)

	if err != nil {
		log.Printf("failed to cast type, %v, %v", err, string(message))
		return nil, err
	}

	return &d, nil
}

func convertToModel(data *model.UserData) (*models.User, error) {
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
