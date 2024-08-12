package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func product(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/product")

	var repo repository.ProductRepoInterface = repository.NewProductRepository(d)
	handler := handlers.NewProduct(repo)

	route.GET("/", handler.GetProducts)
	route.GET("/:uuid", handler.GetProductsDetail)
	route.POST("/", handler.InsertProducts)
	route.DELETE("/:uuid", handler.ProductsDelete)
	route.PATCH("/:uuid", handler.ProductsUpdate)
}
