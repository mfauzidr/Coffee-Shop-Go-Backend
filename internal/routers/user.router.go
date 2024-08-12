package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/middleware"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func user(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/user")

	var repo repository.UserRepoInterface = repository.NewUserRepository(d)
	handler := handlers.NewUserRepository(repo)

	route.GET("/", middleware.AuthMiddleware("admin"), handler.GetUsers)
	route.GET("/:uuid", middleware.AuthMiddleware("admin", "customer"), handler.GetUsersDetail)
	route.POST("/", middleware.AuthMiddleware("admin"), handler.InsertUsers)
	route.PATCH("/:uuid", middleware.AuthMiddleware("admin", "customer"), handler.UsersUpdate)
	route.DELETE("/:uuid", middleware.AuthMiddleware("admin"), handler.UsersDelete)
}
