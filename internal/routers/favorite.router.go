package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/middleware"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func favorite(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/favorite")

	var repo repository.FavoriteRepoInterface = repository.NewFavoriteRepository(d)
	handler := handlers.NewFavoriteRepository(repo)

	route.GET("/", middleware.AuthMiddleware("admin"), handler.GetFavorites)
	route.GET("/:id", middleware.AuthMiddleware("admin", "customer"), handler.GetFavoriteDetail)
	route.POST("/", middleware.AuthMiddleware("customer"), handler.PostFavorite)
	route.DELETE("/:id", middleware.AuthMiddleware("customer"), handler.FavoriteDelete)
	route.PATCH("/:id", middleware.AuthMiddleware("admin", "customer"), handler.PatchFavorite)
}
