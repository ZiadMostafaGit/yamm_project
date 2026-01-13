package handler

import (
	"net/http"
	"yamm-project/app/internal/dto"
	"yamm-project/app/internal/models"
	"yamm-project/app/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	token, err := h.authService.Register(user, req.StoreName)
	if err != nil {

		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusCreated, "Created successfully", gin.H{"token": token})
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	token, err := h.authService.RegisterAdmin(user)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return

	}

	SendResponse(c, http.StatusCreated, "Created successfully", gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.authService.LogIn(req.Email, req.Password)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "Login successfully", gin.H{"token": token})
}
