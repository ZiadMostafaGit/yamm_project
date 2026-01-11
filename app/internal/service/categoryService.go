package service

import (
	"yamm-project/app/internal/models"
	"yamm-project/app/internal/repository"
)

type CategoryService interface {
	CreateCategory(name string) error
	GetAllCategories() ([]models.Category, error)
	UpdateCategory(id uint, category string) error
	DeleteCategoryById(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: repo}
}

func (s *categoryService) CreateCategory(name string) error {
	category := &models.Category{Name: name}
	return s.categoryRepo.Create(category)
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) UpdateCategory(id uint, category string) error {
	return s.categoryRepo.Update(id, category)
}

func (s *categoryService) DeleteCategoryById(id uint) error {
	return s.categoryRepo.Delete(id)
}
