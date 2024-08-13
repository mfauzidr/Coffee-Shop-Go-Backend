package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
)

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := pkg.NewResponse(ctx)
		var header string

		if header = ctx.GetHeader("Authorization"); header == "" {
			response.Unauthorized("Unauthorized", nil)
			ctx.Abort()
			return
		}

		if !strings.Contains(header, "Bearer") {
			response.Unauthorized("Invalid Bearer Token", nil)
			ctx.Abort()
			return
		}

		token := strings.Replace(header, "Bearer ", "", -1)

		check, err := pkg.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Invalid Bearer Token", nil)
			ctx.Abort()
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if check.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			response.Unauthorized("Forbidden: Insufficient permissions", nil)
			ctx.Abort()
			return
		}

		ctx.Set("userUuid", check.UUID)
		ctx.Set("userRole", check.Role)
		ctx.Next()
	}
}
