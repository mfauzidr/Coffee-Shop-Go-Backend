package handlers

import (
	"net/http"
	"strconv"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	*repository.RepoProduct
}

func NewProduct(r *repository.RepoProduct) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) PostProduct(ctx *gin.Context) {
	product := models.Product{}

	if err := ctx.ShouldBind(&product); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.CreateProduct(&product); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
}

func (h *HandlerProduct) GetProducts(ctx *gin.Context) {
	productName := ctx.Query("name")
	minPriceStr := ctx.Query("minPrice")
	maxPriceStr := ctx.Query("maxPrice")
	category := ctx.Query("category")
	sort := ctx.Query("sort")
	pageStr := ctx.Query("page")

	var minPrice, maxPrice, page int
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.Atoi(minPriceStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minPrice"})
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.Atoi(maxPriceStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxPrice"})
			return
		}
	}

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
	} else {
		page = 1
	}

	query := models.ProductQuery{
		Name:     productName,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Category: category,
		Sort:     sort,
		Page:     page,
	}

	data, err := h.GetAllProduct(&query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(*data) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	data, err := h.GetDetailProduct(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching product details"})
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (h *HandlerProduct) ProductDelete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.DeleteProduct(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *HandlerProduct) ProductUpdate(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBind(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	updatedProduct, err := h.UpdateProduct(&product, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}
