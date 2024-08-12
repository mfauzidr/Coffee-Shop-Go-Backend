package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
)

func authRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/auth")

	var authRepo repository.AuthRepoInterface = repository.NewAuthRepository(d)
	handler := handlers.NewAuthHandler(authRepo)

	route.POST("/register", handler.Register)
	route.POST("/login", handler.Login)
}
