//should be deleted

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
)

func (h *AuthHandler) LoginDummy(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Users{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Login failed, please insert email and password", err.Error())
		return
	}

	results, err := h.GetByEmail(body.Email)
	if err != nil {
		response.InternalServerError("Login failed, internal server error", err.Error())
		return
	}

	result := (*results)[0]

	if body.Email != result.Email {
		response.Unauthorized("Wrong email", nil)
		return
	}
	if body.Password != result.Password {
		response.Unauthorized("Wrong password", nil)
		return
	}

	response.Success("Login success", results)
}
