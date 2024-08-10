package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		fmt.Println("Root route accessed")
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to CoffeeShop!"})
	})

	user(router, db)
	product(router, db)
	favorite(router, db)

	return router
}
