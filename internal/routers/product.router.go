package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func product(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/product")

	repo := repository.NewProduct(d)
	handler := handlers.NewProduct(repo)

	route.GET("/", handler.GetProducts)
	route.GET("/:id", handler.GetProductDetail)
	route.POST("/", handler.PostProduct)
	route.DELETE("/:id", handler.ProductDelete)
	route.PATCH("/:id", handler.ProductUpdate)
}
