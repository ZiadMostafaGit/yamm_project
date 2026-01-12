package handler

import (
	"net/http"
	"strconv"
	"yamm-project/app/internal/dto"
	"yamm-project/app/internal/service"

	"github.com/gin-gonic/gin"
)

type FAQHandler struct {
	faqService service.FAQService
}

func NewFAQHandler(fs service.FAQService) *FAQHandler {
	return &FAQHandler{faqService: fs}
}

func (h *FAQHandler) Create(c *gin.Context) {
	var req dto.CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := c.GetString("role")
	var storeIDPtr *uint
	if sid, exists := c.Get("store_id"); exists {
		id := sid.(uint)
		if id != 0 {
			storeIDPtr = &id
		}
	}

	id, err := h.faqService.CreateFAQ(&req, role, storeIDPtr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"faq_id": id})
}

func (h *FAQHandler) GetCustomerView(c *gin.Context) {
	lang := c.Query("lang")

	var storeIDPtr *uint
	if sidStr := c.Query("store_id"); sidStr != "" {
		sid, _ := strconv.ParseUint(sidStr, 10, 64)
		id := uint(sid)
		storeIDPtr = &id
	}

	faqs, err := h.faqService.GetGroupedFAQs(storeIDPtr, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, faqs)
}
func (h *FAQHandler) UpdateVisibility(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateFAQRequest
	newErr := c.ShouldBindJSON(&req)
	if newErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": newErr.Error()})
		return
	}

	role := c.GetString("role")
	otherErr := h.faqService.UpdateFAQVisibility(uint(id), &req, role)
	if otherErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": otherErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visibility updated successfully"})
}

// func (h *FAQHandler) UpdateVisibility(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
// 	var req dto.UpdateFAQRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	role := c.GetString("role")
// 	var sidPtr *uint

// 	sid, exists := c.Get("store_id")
// 	if exists && req.IsGlobal == true {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Cant change to public while it has store"})
// 		return

// 	}
// 	store_id := sid.(uint)
// 	sidPtr = &store_id

// 	if err := h.faqService.UpdateFAQVisibility(uint(id), &req, role, sidPtr); err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Visibility updated"})
// }

func (h *FAQHandler) UpdateTranslations(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req dto.UpdateTranslationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := c.GetString("role")
	var storeIDPtr *uint
	if sid, exists := c.Get("store_id"); exists {
		id := sid.(uint)
		if id != 0 {
			storeIDPtr = &id
		}
	}

	if err := h.faqService.UpdateTranslations(uint(id), &req, role, storeIDPtr); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Translations updated successfully"})
}

func (h *FAQHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	role := c.GetString("role")
	var storeIDPtr *uint
	if sid, exists := c.Get("store_id"); exists {
		id := sid.(uint)
		if id != 0 {
			storeIDPtr = &id
		}
	}

	if err := h.faqService.DeleteFAQ(uint(id), role, storeIDPtr); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FAQ deleted successfully"})
}
