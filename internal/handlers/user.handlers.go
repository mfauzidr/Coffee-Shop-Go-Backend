package handlers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository.UserRepoInterface
}

func NewUserRepository(r repository.UserRepoInterface) *UserHandler {
	return &UserHandler{r}
}

func (h *UserHandler) InsertUsers(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	users := models.Users{}

	if err := ctx.ShouldBind(&users); err != nil {
		response.BadRequest("Create user failed, invalid input", err.Error())
		return
	}

	_, err := govalidator.ValidateStruct(&users)
	if err != nil {
		response.BadRequest("create data failed", err.Error())
		return
	}

	users.Password, err = pkg.HashPassword(users.Password)
	if err != nil {
		response.BadRequest("Register failed", err.Error())
		return
	}

	results, err := h.CreateUser(&users)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_email_key") {
				response.BadRequest("Create user failed, email already exists", err.Error())
				return
			}
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}
	response.Created("User created successfully", results)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "6")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest("Invalid or missing 'page' parameter", err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		response.BadRequest("Invalid or missing 'limit' parameter", err.Error())
		return
	}

	query := models.UsersQuery{
		Page:  page,
		Limit: limit,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest("Invalid query parameter", err.Error())
		return
	}

	data, total, err := h.GetAllUsers(&query)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if len(*data) == 0 {
		response.NotFound("User Not Found", "No users available for the given criteria")
		return
	}

	totalPages := (total + query.Limit - 1) / query.Limit
	meta := &pkg.Meta{
		Total:     total,
		TotalPage: totalPages,
		Page:      query.Page,
		NextPage:  query.Page + 1,
		PrevPage:  query.Page - 1,
	}

	response.GetAllSuccess("User retrieved successfully", data, meta)
}

func (h *UserHandler) GetUsersDetail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	role, _ := ctx.Get("userRole")
	var uuid string

	if role == "customer" {
		if userUuid, ok := ctx.Get("userUuid"); ok {
			uuid = userUuid.(string)
		} else {
			response.Unauthorized("User UUID not found", nil)
			return
		}
	} else {
		uuid = ctx.Param("uuid")
	}

	data, err := h.GetDetailsUser(uuid)
	if err != nil {
		response.NotFound("User Not Found", err.Error())
		return
	}

	response.Success("User details retrieved successfully", data)
}

func (h *UserHandler) UsersUpdate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var input models.Users

	if err := ctx.ShouldBind(&input); err != nil {
		response.BadRequest("Update user failed, invalid input", err.Error())
		return
	}

	_, err := govalidator.ValidateStruct(&input)
	if err != nil {
		response.BadRequest("Update user failed", err.Error())
		return
	}

	input.Password, err = pkg.HashPassword(input.Password)
	if err != nil {
		response.BadRequest("Register failed", err.Error())
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
		if dbTag == "id" {
			continue
		}

		if (fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil()) ||
			(fieldValue.Kind() != reflect.Ptr && fieldValue.Interface() != "") {
			data[dbTag] = fieldValue.Interface()
		}
	}

	role, _ := ctx.Get("userRole")
	var uuid string

	if role == "customer" {
		if userUuid, ok := ctx.Get("userUuid"); ok {
			uuid = userUuid.(string)
		} else {
			response.Unauthorized("User UUID not found", err)
			return
		}
	} else {
		uuid = ctx.Param("uuid")
	}

	updatedUser, err := h.UpdateUser(uuid, data)
	if err != nil {
		response.BadRequest("Update user failed", err.Error())
		return
	}

	response.Success("User updated successfully", updatedUser)
}

func (h *UserHandler) UsersDelete(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	uuid := ctx.Param("uuid")

	data, err := h.DeleteUser(uuid)
	if err != nil {
		if err.Error() == fmt.Sprintf("User with UUID %s not found", uuid) {
			response.NotFound("User Not Found", err.Error())
			return
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	response.Success("User deleted successfully", data)
}
