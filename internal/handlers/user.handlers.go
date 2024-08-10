package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	// "strconv"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type HandlerUsers struct {
	*repository.UsersRepo
}

func NewUsers(r *repository.UsersRepo) *HandlerUsers {
	return &HandlerUsers{r}
}

func (h *HandlerUsers) InsertUsers(ctx *gin.Context) {

	users := models.Users{}

	if err := ctx.ShouldBind(&users); err != nil {
		response := models.Response{
			Status:  "error",
			Message: "invalid input",
			Error:   err.Error(),
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userData := map[string]interface{}{
		"firstName":       users.FirstName,
		"lastName":        users.LastName,
		"phoneNumber":     users.PhoneNumber,
		"address":         users.Address,
		"deliveryAddress": users.DeliveryAddress,
		"image":           users.Image,
		"birthday":        users.Birthday,
		"email":           users.Email,
		"password":        users.Password,
		"role":            users.Role,
		"gender":          users.Gender,
	}

	results, err := h.CreateUser(userData)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_email_key") {
				response := models.Response{
					Status:  "error",
					Message: "Email already exists",
				}
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}
		}
		response := models.Response{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   err.Error(),
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := models.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    results,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *HandlerUsers) GetUsers(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "6")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid or missing 'page' parameter",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid or missing 'limit' parameter",
		})
		return
	}

	query := models.UsersQuery{
		Page:  page,
		Limit: limit,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	data, total, err := h.GetAllUsers(&query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get users",
			Error:   err.Error(),
		})
		return
	}

	if len(*data) == 0 {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "success",
			Message: "No users found",
			Data:    []interface{}{},
		})
		return
	}

	totalPages := (total + query.Limit - 1) / query.Limit
	meta := &models.Meta{
		Total:     total,
		TotalPage: totalPages,
		Page:      query.Page,
		NextPage:  query.Page + 1,
		PrevPage:  query.Page - 1,
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Meta:    meta,
		Data:    data,
	})
}

func (h *HandlerUsers) GetUsersDetail(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := h.GetDetailsUser(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User details retrieved successfully",
		Data:    data,
	})
}

func (h *HandlerUsers) UsersUpdate(ctx *gin.Context) {
	var input models.Users

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(fieldType.Name)
		}

		if (fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil()) ||
			(fieldValue.Kind() != reflect.Ptr && fieldValue.Interface() != "") {
			data[dbTag] = fieldValue.Interface()
		}
	}

	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	fmt.Println(data)

	updatedUser, err := h.UpdateUser(uuid, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Failed to update",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}

func (h *HandlerUsers) UsersDelete(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := h.DeleteUser(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "Failed to Delete. User not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data:    data,
	})
}
