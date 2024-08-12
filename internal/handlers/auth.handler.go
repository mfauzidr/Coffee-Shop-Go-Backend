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
		response.BadRequest("Login failed, insert email or password", err.Error())
		return
	}
	_, err := govalidator.ValidateStruct(&body)
	if err != nil {
		response.BadRequest("Login failed, email is not valid", err.Error())
		return
	}

	result, err := h.GetByEmail(body.Email)
	if err != nil {
		response.BadRequest("Login failed, email is not registered", err.Error())
		return
	}

	err = pkg.VerifyPassword(result.Password, body.Password)
	if err != nil {
		response.Unauthorized("Wrong password", err.Error())
		return
	}

	jwt := pkg.NewJWT(result.UsersUuid, result.Email, result.Role)
	token, err := jwt.GenerateToken()
	if err != nil {
		response.Unauthorized("Failed generate token", err.Error())
		return
	}

	response.Created("login success", token)
}
