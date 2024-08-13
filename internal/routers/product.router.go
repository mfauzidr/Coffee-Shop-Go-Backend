package routers

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/handlers"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/middleware"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func product(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/product")

	var repo repository.ProductRepoInterface = repository.NewProductRepository(d)
	var cld pkg.Cloudinary = *pkg.NewCloudinaryUtil()
	handler := handlers.NewProduct(repo, cld)

	route.GET("/", handler.GetProducts)
	route.GET("/:uuid", handler.GetProductsDetail)
	route.POST("/", middleware.AuthMiddleware("admin"), handler.InsertProducts)
	route.DELETE("/:uuid", middleware.AuthMiddleware("admin"), handler.ProductsDelete)
	route.PATCH("/:uuid", middleware.AuthMiddleware("admin"), handler.ProductsUpdate)
}
