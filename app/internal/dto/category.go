package dto

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
