package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	repository.ProductRepoInterface
}

func NewProduct(r repository.ProductRepoInterface) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) InsertProducts(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	products := models.Product{}

	if err := ctx.ShouldBind(&products); err != nil {
		response.BadRequest("Create product failed, invalid input", err.Error())
		return
	}

	_, err := govalidator.ValidateStruct(&products)
	if err != nil {
		response.BadRequest("Create data failed", err.Error())
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
				response.BadRequest("Create user failed, prodict name already exists", err.Error())
				return
			}
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}
	response.Created("Product created successfully", results)
}

func (h *HandlerProduct) GetProducts(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "6")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest("Invalid or missing 'page' parameter", err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		response.BadRequest("Invalid or missing 'limit' parameter", err.Error())
		return
	}

	query := models.ProductQuery{
		Page:  page,
		Limit: limit,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest("Invalid query parameter", err.Error())
		return
	}

	data, total, err := h.GetAllProduct(&query)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if len(*data) == 0 {
		response.NotFound("User Not Found", "No products available for the given criteria")
		return
	}

	totalPages := (total + query.Limit - 1) / query.Limit
	meta := &pkg.Meta{
		Total:     total,
		TotalPage: totalPages,
		Page:      query.Page,
		NextPage:  query.Page + 1,
		PrevPage:  query.Page - 1,
	}

	response.GetAllSuccess("Product retrieved successfully", data, meta)
}

func (h *HandlerProduct) GetProductsDetail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	uuid := ctx.Param("uuid")

	data, err := h.GetDetailProduct(uuid)
	if err != nil {
		response.NotFound("Product Not Found", err.Error())
		return
	}

	response.Success("Product details retrieved successfully", data)
}

func (h *HandlerProduct) ProductsUpdate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var input models.Product

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := govalidator.ValidateStruct(&input)
	if err != nil {
		response.BadRequest("Create data failed", err.Error())
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
		response.BadRequest("Update product failed, uuid is required", "UUID is Empty")
		return
	}

	fmt.Println(data)

	updatedProduct, err := h.UpdateProduct(uuid, data)
	if err != nil {
		response.BadRequest("Update product failed", err.Error())
		return
	}

	response.Success("Product updated successfully", updatedProduct)
}

func (h *HandlerProduct) ProductsDelete(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	uuid := ctx.Param("uuid")

	data, err := h.DeleteProduct(uuid)
	if err != nil {
		if err.Error() == fmt.Sprintf("Product with UUID %s not found", uuid) {
			response.NotFound("Product Not Found", err.Error())
			return
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	response.Success("User deleted successfully", data)
}
