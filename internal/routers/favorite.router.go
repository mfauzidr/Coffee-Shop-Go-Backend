package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func favorite(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/favorite")

	repo := repository.NewFavorite(d)
	handler := handlers.NewFavorite(repo)

	route.GET("/", handler.GetFavorites)
	route.GET("/:id", handler.GetFavoriteDetail)
	route.POST("/", handler.PostFavorite)
	route.DELETE("/:id", handler.FavoriteDelete)
	route.PATCH("/:id", handler.PatchFavorite)
}
