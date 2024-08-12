package handlers

import (
	"net/http"
	"strconv"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repository.FavoriteRepoInterface
}

func NewFavoriteRepository(r repository.FavoriteRepoInterface) *FavoriteHandler {
	return &FavoriteHandler{r}
}

func (h *FavoriteHandler) PostFavorite(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	favorite := models.PostFavorite{}

	err := ctx.ShouldBind(&favorite)
	if err != nil {
		response.BadRequest("Create failed, invalid input", err.Error())
		return
	}
	results, err := h.CreateFavorite(&favorite)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	response.Created("Create data successfully", results)
}

func (h *FavoriteHandler) GetFavorites(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	data, err := h.GetAllFavorite()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if len(*data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *FavoriteHandler) GetFavoriteDetail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.BadRequest("Failed to retrieve data, invalid input", err.Error())
		return
	}

	data, err := h.GetDetailFavorite(id)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if data == nil {
		response.NotFound("Data not found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *FavoriteHandler) FavoriteDelete(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.BadRequest("Delete data failed, invalid input", err.Error())
		return
	}

	if err := h.DeleteFavorite(id); err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Favorite product deleted successfully"})
}

func (h *FavoriteHandler) PatchFavorite(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var favorite models.UpdateFavorite
	if err := ctx.ShouldBind(&favorite); err != nil {
		response.BadRequest("Update User Favorite Product failed, invalid input", err.Error())
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Update User Favorite Product failed", err.Error())
		return
	}

	data, err := h.UpdateFavorite(id, &favorite)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if data == nil {
		response.NotFound("User Favorite Product Not Found", "No data available for the given criteria")
		return
	}

	response.Success("User deleted successfully", data)
}
