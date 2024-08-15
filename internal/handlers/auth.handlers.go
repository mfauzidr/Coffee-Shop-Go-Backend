package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
)

type AuthHandler struct {
	repository.AuthRepoInterface
}

func NewAuthHandler(authRepo repository.AuthRepoInterface) *AuthHandler {
	return &AuthHandler{authRepo}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Users{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Register failed", err.Error())
		return
	}
	_, err := govalidator.ValidateStruct(&body)
	if err != nil {
		response.BadRequest("Register failed", err.Error())
		return
	}

	body.Password, err = pkg.HashPassword(body.Password)
	if err != nil {
		response.BadRequest("Register failed", err.Error())
		return
	}

	body.Role = "customer"

	result, err := h.RegisterUser(&body)
	if err != nil {
		response.BadRequest("Register failed", err.Error())
		return
	}

	response.Created("Register success", result)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Users{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Login failed, please insert email and password", err.Error())
		return
	}
	_, err := govalidator.ValidateStruct(&body)
	if err != nil {
		response.BadRequest("Login failed, invalid email format", err.Error())
		return
	}

	results, err := h.GetByEmail(body.Email)
	if err != nil {
		response.InternalServerError("Login failed, internal server error", err.Error())
		return
	}

	if results == nil || len(*results) == 0 {
		response.BadRequest("Login failed, email is not registered", "Email not found")
		return
	}

	result := (*results)[0]

	err = pkg.VerifyPassword(result.Password, body.Password)
	if err != nil {
		response.Unauthorized("Wrong password", err.Error())
		return
	}

	jwt := pkg.NewJWT(result.UUID, result.Email, result.Role, result.Id)
	token, err := jwt.GenerateToken()
	if err != nil {
		response.Unauthorized("Failed generate token", err.Error())
		return
	}

	response.Success("login success", token)
}
