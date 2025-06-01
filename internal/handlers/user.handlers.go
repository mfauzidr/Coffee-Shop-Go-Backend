package handlers

import (
	"fmt"
	"math/rand"
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
	pkg.Cloudinary
}

func NewUserRepository(r repository.UserRepoInterface, cld pkg.Cloudinary) *UserHandler {
	return &UserHandler{r, cld}

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
		response.BadRequest("Create User failed", err.Error())
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err == nil {
		mimeType := header.Header.Get("Content-Type")
		if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
			response.BadRequest("Create User failed, upload file failed, file is not supported", nil)
			return
		}

		if header.Size > 2*1024*1024 {
			response.BadRequest("Create User failed, upload file failed, file size exceeds 2 MB", nil)
			return
		}

		randomNumber := rand.Int()
		fileName := fmt.Sprintf("user-image-%d", randomNumber)
		uploadResult, err := h.UploadFile(ctx, file, fileName)
		if err != nil {
			response.BadRequest("Create User failed, upload file failed", err.Error())
			return
		}
		imageURL := uploadResult.SecureURL
		users.Image = &imageURL
	}

	users.Password, err = pkg.HashPassword(users.Password)
	if err != nil {
		response.BadRequest("Create User failed", err.Error())
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
		NextPage:  0,
		PrevPage:  0,
	}

	if query.Page+1 <= totalPages {
		meta.NextPage = query.Page + 1
	}
	if query.Page > 1 {
		meta.PrevPage = query.Page - 1
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

	fmt.Println("Fetching details for user UUID:", uuid)

	data, err := h.GetDetailsUser(uuid)
	if err != nil {
		response.NotFound("User Not Found", err.Error())
		return
	}

	if data == nil {
		response.NotFound("User Not Found", nil)
		return
	}

	response.Success("User details retrieved successfully", data)
}

func (h *UserHandler) UsersUpdate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	input := models.Users{}

	if err := ctx.ShouldBind(&input); err != nil {
		response.BadRequest("Update user failed, invalid input", err.Error())
		return
	}

	// Validasi minimal
	if input.Email != "" && !govalidator.IsEmail(input.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}
	if input.Password != "" && len(input.Password) < 6 {
		response.BadRequest("Password minimal 6 karakter", nil)
		return
	}
	if input.PhoneNumber != nil && *input.PhoneNumber != "" && !govalidator.IsNumeric(*input.PhoneNumber) {
		response.BadRequest("Nomor telepon hanya boleh angka", nil)
		return
	}

	// Upload file jika ada
	file, header, err := ctx.Request.FormFile("image")
	if err == nil {
		mimeType := header.Header.Get("Content-Type")
		if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
			response.BadRequest("Upload failed: file type not supported", nil)
			return
		}
		if header.Size > 2*1024*1024 {
			response.BadRequest("Upload failed: file size exceeds 2MB", nil)
			return
		}
		fileName := fmt.Sprintf("user-image-%d", rand.Int())
		uploadResult, err := h.UploadFile(ctx, file, fileName)
		if err != nil {
			response.BadRequest("Upload failed", err.Error())
			return
		}
		input.Image = &uploadResult.SecureURL
	}

	// Hash password jika diberikan
	if input.Password != "" {
		hashedPassword, err := pkg.HashPassword(input.Password)
		if err != nil {
			response.BadRequest("Hashing failed", err.Error())
			return
		}
		input.Password = hashedPassword
	}

	// UUID dari context atau param
	var uuid string
	role, _ := ctx.Get("userRole")
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

	updatedUser, err := h.UpdateUser(uuid, &input)
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
