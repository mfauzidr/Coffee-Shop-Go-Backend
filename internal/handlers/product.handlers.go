package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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

func (h *HandlerProduct) InsertProducts(ctx *gin.Context) {
	products := models.Product{}

	if err := ctx.ShouldBind(&products); err != nil {
		response := models.Response{
			Status:  "error",
			Message: "invalid input",
			Error:   err.Error(),
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	productData := map[string]interface{}{
		"name":        products.Name,
		"description": products.Description,
		"price":       products.Price,
		"category":    products.Category,
		"image":       products.Image,
	}

	results, err := h.CreateProduct(productData)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "unique_name") {
				response := models.Response{
					Status:  "error",
					Message: "Products already exists",
				}
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}
		}
		response := models.Response{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   err.Error(),
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := models.Response{
		Status:  "success",
		Message: "Product created successfully",
		Data:    results,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *HandlerProduct) GetProducts(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "6")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid or missing 'page' parameter",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid or missing 'limit' parameter",
		})
		return
	}

	query := models.ProductQuery{
		Page:  page,
		Limit: limit,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	data, total, err := h.GetAllProduct(&query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get products",
			Error:   err.Error(),
		})
		return
	}

	if len(*data) == 0 {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "success",
			Message: "No products found",
			Data:    []interface{}{},
		})
		return
	}

	totalPages := (total + query.Limit - 1) / query.Limit
	meta := &models.Meta{
		Total:     total,
		TotalPage: totalPages,
		Page:      query.Page,
		NextPage:  query.Page + 1,
		PrevPage:  query.Page - 1,
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Products retrieved successfully",
		Meta:    meta,
		Data:    data,
	})
}

func (h *HandlerProduct) GetProductsDetail(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := h.GetDetailProduct(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Product details retrieved successfully",
		Data:    data,
	})
}

func (h *HandlerProduct) ProductsUpdate(ctx *gin.Context) {
	var input models.Product

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(fieldType.Name)
		}
		if dbTag == "id" {
			continue
		}

		if (fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil()) ||
			(fieldValue.Kind() != reflect.Ptr && fieldValue.Interface() != "") {
			data[dbTag] = fieldValue.Interface()
		}
	}

	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductUuid is required"})
		return
	}

	fmt.Println(data)

	upatedProduct, err := h.UpdateProduct(uuid, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Failed to update",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Product updated successfully",
		Data:    upatedProduct,
	})
}

func (h *HandlerProduct) ProductsDelete(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := h.DeleteProduct(uuid)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with UUID %s not found", uuid) {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "error",
				Message: "product not found",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "product deleted successfully",
		Data:    data,
	})
}
