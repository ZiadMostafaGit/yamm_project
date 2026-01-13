package handler

import (
	"net/http"
	"strconv"
	"yamm-project/app/internal/dto"
	"yamm-project/app/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(cs service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

func (ch *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	otherErr := ch.categoryService.CreateCategory(req.Name)
	if otherErr != nil {
		SendResponse(c, http.StatusBadRequest, otherErr.Error(), nil)
		return
	}

	SendResponse(c, http.StatusCreated, "Category created successfully", nil)
}

func (ch *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		SendResponse(c, http.StatusBadRequest, "No id found", nil)
		return
	}

	newId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, "Invalid id number", nil)
		return
	}

	newErr := ch.categoryService.DeleteCategoryById(uint(newId))
	if newErr != nil {
		SendResponse(c, http.StatusBadRequest, "operation failed", nil)
		return
	}

	SendResponse(c, http.StatusOK, "category deleted successfully", nil)
}

func (ch *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := ch.categoryService.GetAllCategories()
	if err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "Categories fetched successfully", categories)
}

func (ch *CategoryHandler) Update(c *gin.Context) {
	newcategoryId := c.Param("id")
	Id, err := strconv.ParseInt(newcategoryId, 10, 64)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, "Invalid Id", nil)
		return
	}

	var req dto.CategoryUpdateRequest
	newErr := c.ShouldBindJSON(&req)
	if newErr != nil {
		SendResponse(c, http.StatusBadRequest, newErr.Error(), nil)
		return
	}

	theardErr := ch.categoryService.UpdateCategory(uint(Id), req.Name)
	if theardErr != nil {
		SendResponse(c, http.StatusBadRequest, theardErr.Error(), nil)
		return
	}

	SendResponse(c, http.StatusAccepted, "Category updated successfully", nil)
}
