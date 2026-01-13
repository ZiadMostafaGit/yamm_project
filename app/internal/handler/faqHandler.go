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
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
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
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusCreated, "FAQ created successfully", gin.H{"faq_id": id})
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
		SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "FAQs fetched successfully", faqs)
}

func (h *FAQHandler) UpdateVisibility(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.UpdateFAQRequest
	newErr := c.ShouldBindJSON(&req)
	if newErr != nil {
		SendResponse(c, http.StatusBadRequest, newErr.Error(), nil)
		return
	}

	role := c.GetString("role")
	otherErr := h.faqService.UpdateFAQVisibility(uint(id), &req, role)
	if otherErr != nil {
		SendResponse(c, http.StatusBadRequest, otherErr.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "Visibility updated successfully", nil)
}

func (h *FAQHandler) UpdateTranslations(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req dto.UpdateTranslationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
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
		SendResponse(c, http.StatusForbidden, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "Translations updated successfully", nil)
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
		SendResponse(c, http.StatusForbidden, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "FAQ deleted successfully", nil)
}
