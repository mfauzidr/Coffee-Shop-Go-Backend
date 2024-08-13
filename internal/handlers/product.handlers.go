package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
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
	pkg.Cloudinary
}

func NewProduct(r repository.ProductRepoInterface, cld pkg.Cloudinary) *HandlerProduct {
	return &HandlerProduct{r, cld}
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

	file, header, err := ctx.Request.FormFile("image")
	if err == nil {
		mimeType := header.Header.Get("Content-Type")
		if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
			response.BadRequest("Create User failed, upload file failed, file is not supported", nil)
			return
		}

		if header.Size > 2*1024*1024 {
			response.BadRequest("Create User failed, upload file failed, file size exceeds 2 MB", nil)
			return
		}

		randomNumber := rand.Int()
		fileName := fmt.Sprintf("user-image-%d", randomNumber)
		uploadResult, err := h.UploadFile(ctx, file, fileName)
		if err != nil {
			response.BadRequest("Create User failed, upload file failed", err.Error())
			return
		}
		imageURL := uploadResult.SecureURL
		products.Image = &imageURL
	}

	results, err := h.CreateProduct(&products)
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

	file, header, err := ctx.Request.FormFile("image")

	if err == nil {
		mimeType := header.Header.Get("Content-Type")
		fmt.Println(mimeType)
		if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
			response.BadRequest("Update User failed, upload file failed, file is not supported", nil)
			return
		}

		if header.Size > 2*1024*1024 {
			response.BadRequest("Update User failed, upload file failed, file size exceeds 2 MB", nil)
			return
		}

		randomNumber := rand.Int()
		fileName := fmt.Sprintf("user-image-%d", randomNumber)
		uploadResult, err := h.UploadFile(ctx, file, fileName)
		if err != nil {
			response.BadRequest("Update User failed, upload file failed", err.Error())
			return
		}
		imageURL := uploadResult.SecureURL
		input.Image = &imageURL
	}

	uuid := ctx.Param("uuid")
	if uuid == "" {
		response.BadRequest("Update product failed, uuid is required", "UUID is Empty")
		return
	}

	updatedProduct, err := h.UpdateProduct(uuid, &input)
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
