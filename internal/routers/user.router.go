package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func user(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/user")

	var repo repository.UserRepoInterface = repository.NewUserRepository(d)
	handler := handlers.NewUserRepository(repo)

	route.GET("/", handler.GetUsers)
	route.GET("/:uuid", handler.GetUsersDetail)
	route.POST("/", handler.InsertUsers)
	route.PATCH("/:uuid", handler.UsersUpdate)
	route.DELETE("/:uuid", handler.UsersDelete)
}
