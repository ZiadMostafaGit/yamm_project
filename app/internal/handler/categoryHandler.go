package handler

import (
	"net/http"
	"strconv"
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
	var req CategoryCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otherErr := ch.categoryService.CreateCategory(req.Name)
	if otherErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": otherErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})

}

func (ch *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "No id found"})
		return
	}
	newId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id number"})
		return
	}
	newErr := ch.categoryService.DeleteCategoryById(uint(newId))
	if newErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "operation faild"})
		return

	}

	c.JSON(http.StatusCreated, gin.H{"message": "category deleted successfully"})
}

func (ch *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := ch.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})

}

func (ch *CategoryHandler) Update(c *gin.Context) {
	newcategoryId := c.Param("id")
	Id, err := strconv.ParseInt(newcategoryId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return

	}

	var req CategoryUpdateRequest
	newErr := c.ShouldBindJSON(&req)
	if newErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": newErr.Error()})
		return

	}

	theardErr := ch.categoryService.UpdateCategory(uint(Id), req.Name)
	if theardErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": theardErr.Error()})
		return

	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Category updated successfully"})

}
